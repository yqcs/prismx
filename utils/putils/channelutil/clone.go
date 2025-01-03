package channelutil

import (
	"context"
	"log"
	"sync"

	errorutil "prismx_cli/utils/putils/errors"
)

// CloneOptions provides options for Cloning channels
type CloneOptions struct {
	MaxDrain  int // Max buffers to drain at once(default 3)
	Threshold int // Threshold(default 5) is buffer length at which drains are activated
}

// CloneChannels provides method to Clone channels
type CloneChannels[T any] struct {
	// Options
	opts *CloneOptions
	Log  *log.Logger

	// Internal
	wg sync.WaitGroup
}

// NewCloneChannels returns new instance of CloneChannels
func NewCloneChannels[T any](opts *CloneOptions) *CloneChannels[T] {
	if opts == nil {
		opts = &CloneOptions{}
	}
	if opts.MaxDrain == 0 {
		opts.MaxDrain = 3
	}
	if opts.Threshold == 0 {
		opts.Threshold = 5
	}
	return &CloneChannels[T]{
		opts: opts,
	}
}

// Clone takes data from source channel(src) and sends them to sinks(send only channel) without being totally unfair
func (s *CloneChannels[T]) Clone(ctx context.Context, src chan T, sinks ...chan<- T) error {
	if src == nil {
		return errorutil.New("source channel is nil").WithTag("Clone", "channel")
	}

	// check if all sinks are not nil
	for _, ch := range sinks {
		if ch == nil {
			return errorutil.New("nil sink found").WithTag("Clone", "channel")
		}
	}

	// Worker Only Supports 5 sinks for now
	if len(sinks)%5 != 0 {
		remaining := 5 - (len(sinks) % 5)
		for i := 0; i < remaining; i++ {
			// add nil channels, these are automatically kicked out of select
			sinks = append(sinks, nil)
		}
	}

	if len(sinks) == 5 {
		s.wg.Add(1)
		go s.cloneChanWorker(ctx, src, sinks...)
		return nil
	}

	/*
		If sinks > 5
		relay channels are used that relay data from root node to leaf node (i.e in this case channel)

		1. sinks are grouped into 5 with 1 relay channel for each group
		2. Each group is passed to worker
		3. Relay are fed to Clone i.e Recursion
	*/
	/*
			Ex:
					   $ 			 <-  Source Channel
				     /   \
				    $  	  $			 <-  Relay Channels
			       / \ 	 / \
			      $   $ $   $		 <-  Leaf Channels (i.e Sinks)

		*Simplicity purpose 2 childs are shown for each nodebut each node(except root node) has 5 childs
	*/

	// create groups of 5 sinks
	groups := [][]chan<- T{}
	tmp := []chan<- T{}
	for i, v := range sinks {
		if i != 0 && i%5 == 0 {
			groups = append(groups, tmp)
			tmp = []chan<- T{}
		}
		tmp = append(tmp, v)
	}
	if len(tmp) > 0 {
		groups = append(groups, tmp)
	}

	// for each group create relay channel
	relaychannels := []chan<- T{}
	// launch worker groups
	for _, v := range groups {
		relay := make(chan T)
		relaychannels = append(relaychannels, relay)
		s.wg.Add(1)
		go s.cloneChanWorker(ctx, relay, v...)
	}

	// recursion use sources to feed relays
	return s.Clone(ctx, src, relaychannels...)
}

// CloneChanWorker is actual worker goroutine
func (s *CloneChannels[T]) cloneChanWorker(ctx context.Context, src chan T, sinkchans ...chan<- T) {
	defer func() {
		for _, v := range sinkchans {
			if v != nil {
				close(v)
			}
		}
		s.wg.Done()
	}()
	if src == nil {
		if s.Log != nil {
			s.Log.Println("Error: source channel is nil")
		}
		return
	}
	if len(sinkchans) != 5 {
		if s.Log != nil {
			s.Log.Printf("Error: expected total sinks 5 but got %v", len(sinkchans))
		}
		return
	}

	sink := map[int]chan<- T{}
	count := 0
	for _, v := range sinkchans {
		sink[count] = v
		count++
	}

	// backlog tracks sink channels whose buffers have reached threshold
	backlog := map[int]struct{}{}
	buffer := map[int][]T{}

	// Helper Functions
	// addToBuff adds data to buffers where data was not sent
	// since channel was not available at that time
	addToBuff := func(id int, value T) {
		for sid, ch := range sink {
			if ch != nil && id != sid {
				if buffer[sid] == nil {
					buffer[sid] = []T{}
				}
				// add to buffer
				buffer[sid] = append(buffer[sid], value)
				if len(buffer[sid]) == s.opts.Threshold {
					backlog[sid] = struct{}{}
				}
			}
		}
	}

	//drain buffer of given channel
	drainAndReset := func() {
		// drain buffer of given channel since threshold has been breached
		// get pseudo random channel using map
		count := 0
		for chanID := range backlog {
			if sink[chanID] != nil {
				for _, item := range buffer[chanID] {
					select {
					case <-ctx.Done():
						return
					case sink[chanID] <- item:
					}
				}
				buffer[chanID] = []T{}
				delete(backlog, chanID)
				count++
				if count == s.opts.MaxDrain {
					// skip for now
					return
				}
			}
		}
	}

	// this is main loop of worker where source channels sends data to whatever sink is available at that time
	// and buffers data for channels that are not available at that time. If buffer of any channel reaches threshold
	// then blocking operation of sending data to that sink channel is performed
forloop:
	for {
		switch {
		case len(backlog) > 0:
			// if buffer of any channel has reached threshold
			drainAndReset()
			// if it is true
		default:
			// send to sinks
			w, ok := <-src
			if !ok {
				break forloop
			}
			select {
			case <-ctx.Done():
				return
			case sink[0] <- w:
				addToBuff(0, w)
			case sink[1] <- w:
				addToBuff(1, w)
			case sink[2] <- w:
				addToBuff(2, w)
			case sink[3] <- w:
				addToBuff(3, w)
			case sink[4] <- w:
				addToBuff(4, w)
			}
		}
	}

	// commit all buffers when source channel is closed
	for id, ch := range sink {
		if ch != nil {
			if len(buffer[id]) > 0 {
				for _, item := range buffer[id] {
					select {
					case <-ctx.Done():
						return
					case ch <- item:
					}
				}
			}
		}
	}
}

// Waits until cloning is finished
func (s *CloneChannels[T]) Wait() {
	s.wg.Wait()
}
