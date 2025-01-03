package batcher

import (
	"time"
)

// FlushCallback is the callback function that will be called when the batcher is full or the flush interval is reached
type FlushCallback[T any] func([]T)

// Batcher is a batcher for any type of data
type Batcher[T any] struct {
	maxCapacity   int
	flushInterval *time.Duration
	flushCallback FlushCallback[T]

	incomingData chan T
	full         chan bool
	mustExit     chan bool
	done         chan bool
}

// BatcherOption is the option for the batcher
type BatcherOption[T any] func(*Batcher[T])

// WithMaxCapacity sets the max capacity of the batcher
func WithMaxCapacity[T any](maxCapacity int) BatcherOption[T] {
	return func(b *Batcher[T]) {
		b.maxCapacity = maxCapacity
	}
}

// WithFlushInterval sets the optional flush interval of the batcher
func WithFlushInterval[T any](flushInterval time.Duration) BatcherOption[T] {
	return func(b *Batcher[T]) {
		b.flushInterval = &flushInterval
	}
}

// WithFlushCallback sets the flush callback of the batcher
func WithFlushCallback[T any](fn FlushCallback[T]) BatcherOption[T] {
	return func(b *Batcher[T]) {
		b.flushCallback = fn
	}
}

// New creates a new batcher
func New[T any](opts ...BatcherOption[T]) *Batcher[T] {
	batcher := &Batcher[T]{
		full:     make(chan bool),
		mustExit: make(chan bool, 1),
		done:     make(chan bool, 1),
	}
	for _, opt := range opts {
		opt(batcher)
	}
	batcher.incomingData = make(chan T, batcher.maxCapacity)
	if batcher.flushCallback == nil {
		panic("batcher: flush callback is required")
	}
	if batcher.maxCapacity <= 0 {
		panic("batcher: max capacity must be greater than 0")
	}
	return batcher
}

// Append appends data to the batcher
func (b *Batcher[T]) Append(d ...T) {
	for _, item := range d {
		if !b.put(item) {
			// will wait until space available
			b.full <- true
			b.incomingData <- item
		}
	}
}

func (b *Batcher[T]) put(d T) bool {
	// try to append the data
	select {
	case b.incomingData <- d:
		return true
	default:
		// channel is full
		return false
	}
}

func (b *Batcher[T]) run() {
	// consume all items in the queue
	defer func() {
		b.doCallback()
		close(b.done)
	}()

	var timer *time.Timer
	var flushInterval time.Duration
	if b.flushInterval != nil {
		flushInterval = *b.flushInterval
		timer = time.NewTimer(flushInterval)
		b.runWithTimer(timer, flushInterval)
		return
	}
	b.runWithoutTimer()
}

// runWithTimer runs the batcher with timer
func (b *Batcher[T]) runWithTimer(timer *time.Timer, flushInterval time.Duration) {
	for {
		select {
		case <-timer.C:
			b.doCallback()
			timer.Reset(flushInterval)
		case <-b.full:
			if !timer.Stop() {
				<-timer.C
			}
			b.doCallback()
			timer.Reset(flushInterval)
		case <-b.mustExit:
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}
}

// runWithoutTimer runs the batcher without timer
func (b *Batcher[T]) runWithoutTimer() {
	for {
		select {
		case <-b.full:
			b.doCallback()
		case <-b.mustExit:
			return
		}
	}
}

func (b *Batcher[T]) doCallback() {
	n := len(b.incomingData)
	if n == 0 {
		return
	}
	items := make([]T, n)

	var k int
	for item := range b.incomingData {
		items[k] = item
		k++
		if k >= n {
			break
		}
	}
	b.flushCallback(items)
}

// Run starts the batcher
func (b *Batcher[T]) Run() {
	go b.run()
}

// Stop stops the batcher
func (b *Batcher[T]) Stop() {
	b.mustExit <- true
}

// WaitDone waits until the batcher is done
func (b *Batcher[T]) WaitDone() {
	<-b.done
}
