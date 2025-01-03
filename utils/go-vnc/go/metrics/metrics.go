package metrics

import (
	"log"
	"math"
)

type Metric interface {
	// Adjust increments or decrements the metric value.
	Adjust(int64)

	// Increment increases the metric value by one.
	Increment()

	// Name returns the varz name.
	Name() string

	// Reset the metric.
	Reset()

	// Value returns the current metric value.
	Value() uint64
}

// Counter provides a simple monotonically incrementing counter.
type Counter struct {
	name string
	val  uint64
}

func (c *Counter) Adjust(val int64) {
	log.Fatal("A Counter metric cannot be adjusted")
}

func (c *Counter) Increment() {
	c.val++
}

func (c *Counter) Name() string {
	return c.name
}

func (c *Counter) Reset() {
	c.val = 0
}

func (c *Counter) Value() uint64 {
	return c.val
}

// The Gauge type represents a non-negative integer, which may increase or
// decrease, but shall never exceed the maximum value.
type Gauge struct {
	name string
	val  uint64
}

// Adjust allows one to increase or decrease a metric.
func (g *Gauge) Adjust(val int64) {
	// The value is positive.
	if val > 0 {
		if g.val == math.MaxUint64 {
			return
		}
		v := g.val + uint64(val)
		if v > g.val {
			g.val = v
			return
		}
		// The value wrapped, so set to maximum allowed value.
		g.val = math.MaxUint64
		return
	}

	// The value is negative.
	v := g.val - uint64(-val)
	if v < g.val {
		g.val = v
		return
	}
	// The value wrapped, so set to zero.
	g.val = 0
}

func (g *Gauge) Increment() {
	log.Fatal("A Gauge metric cannot be adjusted")
}

func (g *Gauge) Name() string {
	return g.name
}

func (g *Gauge) Reset() {
	g.val = 0
}

func (g *Gauge) Value() uint64 {
	return g.val
}
