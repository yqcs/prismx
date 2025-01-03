package channelutil_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"prismx_cli/utils/putils/channelutil"
)

var debug bool = false

// Helper Functions
var FillSource = func(work chan<- struct{}, count int, sg *sync.WaitGroup) {
	defer sg.Done()
	// send work
	for i := 100; i < 100+count; i++ {
		if debug {
			log.Println("sending work")
		}
		work <- struct{}{}
	}
	close(work)
}

var DrainSinkChan = func(ch <-chan struct{}, sleeptime time.Duration, chanId string, wg *sync.WaitGroup) {
	defer wg.Done()
	if ch == nil {
		fmt.Println("worker chan is nil")
	}
	for {
		val, ok := <-ch
		if !ok {
			if debug {
				log.Printf("sink channel %v closed\n", chanId)
			}
			break
		}
		// time.Sleep(sleeptime) // to simulate blocking network i/o
		if debug {
			log.Printf("completed work %v at chan %v\n", val, chanId)
		}
	}
}

func TestCloneCounter(t *testing.T) {
	// test Clone count
	for n := 1; n < 101; n++ {
		source := make(chan struct{})
		sinks := []chan<- struct{}{}
		for i := 0; i < n; i++ {
			sinks = append(sinks, make(chan struct{}))
		}
		cloneOpts := &channelutil.CloneOptions{
			MaxDrain:  3,
			Threshold: 5,
		}
		cloner := channelutil.NewCloneChannels[struct{}](cloneOpts)
		err := cloner.Clone(context.TODO(), source, sinks...)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestJoinCounter(t *testing.T) {
	// test Clone count
	for n := 1; n < 101; n++ {
		sink := make(chan struct{})
		sources := []<-chan struct{}{}
		for i := 0; i < n; i++ {
			sources = append(sources, make(chan struct{}))
		}
		jchan := channelutil.JoinChannels[struct{}]{}
		err := jchan.Join(context.TODO(), sink, sources...)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestCloneN(t *testing.T) {
	for n := 2; n < 44; n++ {
		source := make(chan struct{})
		sinks := []chan struct{}{}
		wg := &sync.WaitGroup{}
		sg := &sync.WaitGroup{}
		for i := 0; i < n; i++ {
			sinks = append(sinks, make(chan struct{}))
			wg.Add(1)
			go DrainSinkChan(sinks[i], time.Millisecond, strconv.Itoa(i), wg)
		}

		// Clone channel
		cloneOpts := &channelutil.CloneOptions{
			MaxDrain:  3,
			Threshold: 5,
		}
		recOnly := []chan<- struct{}{}
		for _, v := range sinks {
			recOnly = append(recOnly, v)
		}
		cloner := channelutil.NewCloneChannels[struct{}](cloneOpts)

		err := cloner.Clone(context.TODO(), source, recOnly...)
		if err != nil {
			t.Error(err)
		}

		sg.Add(1)
		FillSource(source, 10, sg)
		sg.Wait()
		wg.Wait()
	}
}

func TestJoinN(t *testing.T) {
	for N := 2; N < 44; N++ {
		sink := make(chan struct{})
		sg := &sync.WaitGroup{}
		// Drain Sink
		sg.Add(1)
		go DrainSinkChan(sink, time.Millisecond, "drain", sg)
		sources := []chan struct{}{}
		for i := 0; i < N; i++ {
			sources = append(sources, make(chan struct{}))
		}
		// create joiner
		joinch := channelutil.JoinChannels[struct{}]{}
		sndOnly := []<-chan struct{}{}
		for _, v := range sources {
			sndOnly = append(sndOnly, v)
		}
		err := joinch.Join(context.Background(), sink, sndOnly...)
		if err != nil {
			t.Error(err)
		}

		srcgrp := &sync.WaitGroup{}
		// Now start sending data to sources
		for _, v := range sources {
			srcgrp.Add(1)
			go FillSource(v, 10, srcgrp)
		}
		srcgrp.Wait()
		sg.Wait()
	}
}

func TestNIntegration(t *testing.T) {
	// Integration Test i.e Clone and join

	domainChan := make(chan struct{})
	sources := channelutil.CreateNChannels[struct{}](100, 0)
	task := make(chan struct{})

	controller := &sync.WaitGroup{}
	controller.Add(2)

	// send data to domainChan
	go FillSource(domainChan, 100, controller)
	// recieve tasks from all sources
	go DrainSinkChan(task, time.Millisecond, "drain", controller)

	// Clone all channels
	cloneOpts := &channelutil.CloneOptions{
		MaxDrain:  3,
		Threshold: 5,
	}
	srcs := []chan<- struct{}{}
	for _, v := range sources {
		srcs = append(srcs, v)
	}
	cloner := channelutil.NewCloneChannels[struct{}](cloneOpts)
	err := cloner.Clone(context.TODO(), domainChan, srcs...)
	if err != nil {
		t.Error(err)
	}

	// join all channels
	jchan := channelutil.JoinChannels[struct{}]{}
	jsrcs := []<-chan struct{}{}
	for _, v := range sources {
		jsrcs = append(jsrcs, v)
	}
	err = jchan.Join(context.TODO(), task, jsrcs...)
	if err != nil {
		t.Error(err)
	}
	controller.Wait()
}
