package semaphore

import (
	"context"
	"errors"
	"math"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

type Semaphore struct {
	sem         *semaphore.Weighted
	initialSize atomic.Int64
	maxSize     atomic.Int64
	currentSize atomic.Int64
}

func New(size int64) (*Semaphore, error) {
	maxSize := int64(math.MaxInt64)
	s := &Semaphore{
		sem: semaphore.NewWeighted(maxSize),
	}
	s.initialSize.Store(size)
	s.maxSize.Store(maxSize)
	s.currentSize.Store(size)
	err := s.sem.Acquire(context.Background(), s.maxSize.Load()-s.initialSize.Load())
	return s, err
}

func (s *Semaphore) Acquire(ctx context.Context, n int64) error {
	return s.sem.Acquire(ctx, n)
}

func (s *Semaphore) Release(n int64) {
	s.sem.Release(n)
}

// Vary capacity by x - it's internally enqueued as a normal Acquire/Release operation as other Get/Put
// but tokens are held internally
func (s *Semaphore) Vary(ctx context.Context, x int64) error {
	switch {
	case x > 0:
		s.sem.Release(x)
		s.currentSize.Add(x)
		return nil
	case x < 0:
		err := s.sem.Acquire(ctx, x)
		if err != nil {
			return err
		}
		s.currentSize.Add(x)
		return nil
	default:
		return errors.New("x is zero")
	}
}

// Current size of the semaphore
func (s *Semaphore) Size() int64 {
	return s.currentSize.Load()
}

// Nominal size of the sempahore
func (s *Semaphore) InitialSize() int64 {
	return s.initialSize.Load()
}
