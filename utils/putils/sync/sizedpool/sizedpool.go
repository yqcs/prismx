package sizedpool

import (
	"context"
	"errors"
	"sync"

	"prismx_cli/utils/putils/sync/semaphore"
)

type PoolOption[T any] func(*SizedPool[T]) error

func WithSize[T any](size int64) PoolOption[T] {
	return func(sz *SizedPool[T]) error {
		if size <= 0 {
			return errors.New("size must be positive")
		}
		var err error
		sz.sem, err = semaphore.New(size)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithPool[T any](p *sync.Pool) PoolOption[T] {
	return func(sz *SizedPool[T]) error {
		sz.pool = p
		return nil
	}
}

type SizedPool[T any] struct {
	sem  *semaphore.Semaphore
	pool *sync.Pool
}

func New[T any](options ...PoolOption[T]) (*SizedPool[T], error) {
	sz := &SizedPool[T]{}
	for _, option := range options {
		if err := option(sz); err != nil {
			return nil, err
		}
	}
	return sz, nil
}

func (sz *SizedPool[T]) Get(ctx context.Context) (T, error) {
	if sz.sem != nil {
		if err := sz.sem.Acquire(ctx, 1); err != nil {
			var t T
			return t, err
		}
	}
	return sz.pool.Get().(T), nil
}

func (sz *SizedPool[T]) Put(x T) {
	if sz.sem != nil {
		sz.sem.Release(1)
	}
	sz.pool.Put(x)
}

// Vary capacity by x - it's internally enqueued as a normal Acquire/Release operation as other Get/Put
// but tokens are held internally
func (sz *SizedPool[T]) Vary(ctx context.Context, x int64) error {
	return sz.sem.Vary(ctx, x)
}

// Current size of the pool
func (sz *SizedPool[T]) Size() int64 {
	return sz.sem.Size()
}
