package trace

import (
	"errors"
	"math"
	"runtime"
	"sync"
	"time"

	"prismx_cli/utils/putils/generic"
)

const (
	// DefaultMemorySnapshotInterval is the default interval for taking memory snapshots
	DefaultMemorySnapshotInterval = 100 * time.Millisecond
)

type MemorySnapshot struct {
	Time  time.Time
	Alloc uint64
}

type Metrics struct {
	StartTime         time.Time
	FinishTime        time.Time
	ExecutionDuration time.Duration
	Snapshots         []MemorySnapshot
	MinAllocMemory    uint64
	MaxAllocMemory    uint64
	AvgAllocMemory    uint64
}

type FunctionContext struct {
	strategy ActionStrategy
	action   func()
	before   func()
	after    func()
}

func (f *FunctionContext) Execute() {
	if f.before != nil {
		f.before()
	}

	f.strategy.Before()
	f.action()
	f.strategy.After()

	if f.after != nil {
		f.after()
	}
}

type ActionStrategy interface {
	Before()
	After()
	GetMetrics() *Metrics
}

type DefaultStrategy struct {
	metrics generic.Lockable[*Metrics]
	ticker  *time.Ticker
	done    chan bool
	wg      sync.WaitGroup
}

func (d *DefaultStrategy) Before() {
	d.metrics.Do(func(m *Metrics) {
		m.StartTime = time.Now()

		d.ticker = time.NewTicker(DefaultMemorySnapshotInterval)
		d.done = make(chan bool)
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for {
				select {
				case <-d.done:
					return
				case t := <-d.ticker.C:
					var mem runtime.MemStats
					runtime.ReadMemStats(&mem)
					m.Snapshots = append(m.Snapshots, MemorySnapshot{
						Time:  t,
						Alloc: mem.Alloc,
					})
				}
			}
		}()
	})
}

func (d *DefaultStrategy) After() {
	close(d.done)
	d.wg.Wait()
	d.ticker.Stop()
	d.metrics.Do(func(m *Metrics) {
		m.FinishTime = time.Now()
		m.ExecutionDuration = m.FinishTime.Sub(m.StartTime)

		var totalMemory uint64 = 0
		if len(m.Snapshots) > 0 {
			m.MinAllocMemory = m.Snapshots[0].Alloc
			m.MaxAllocMemory = m.Snapshots[0].Alloc

			for _, s := range m.Snapshots {
				if s.Alloc < m.MinAllocMemory {
					m.MinAllocMemory = s.Alloc
				}
				m.MinAllocMemory = uint64(math.Min(float64(m.MinAllocMemory), float64(s.Alloc)))
				m.MaxAllocMemory = uint64(math.Max(float64(m.MaxAllocMemory), float64(s.Alloc)))
				totalMemory += s.Alloc
			}
			m.AvgAllocMemory = totalMemory / uint64(len(m.Snapshots))
		}
	})

}

func (d *DefaultStrategy) GetMetrics() *Metrics {
	var metrics *Metrics
	d.metrics.Do(func(m *Metrics) {
		metrics = m
	})
	return metrics
}

type TraceOptions struct {
	strategy ActionStrategy
	before   func()
	after    func()
}

type TraceOptionSetter func(opts *TraceOptions)

func WithStrategy(s ActionStrategy) TraceOptionSetter {
	return func(opts *TraceOptions) {
		opts.strategy = s
	}
}

func WithBefore(b func()) TraceOptionSetter {
	return func(opts *TraceOptions) {
		opts.before = b
	}
}

func WithAfter(a func()) TraceOptionSetter {
	return func(opts *TraceOptions) {
		opts.after = a
	}
}

func Trace(f func(), setters ...TraceOptionSetter) (*Metrics, error) {
	opts := &TraceOptions{
		strategy: &DefaultStrategy{metrics: generic.Lockable[*Metrics]{V: &Metrics{}}},
	}

	// Apply option if provided
	for _, setter := range setters {
		setter(opts)
	}

	if opts.strategy == nil {
		return nil, errors.New("strategy should not be nil")
	}

	context := &FunctionContext{
		strategy: opts.strategy,
		action:   f,
		before:   opts.before,
		after:    opts.after,
	}

	context.Execute()
	return opts.strategy.GetMetrics(), nil
}
