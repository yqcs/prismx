package async

import "context"

// Future mimics the async/await paradigm
type Future[T any] interface {
	Await() (T, error)
}

type future[T any] struct {
	await func(ctx context.Context) (T, error)
}

func (f future[T]) Await() (T, error) {
	return f.await(context.Background())
}

func Exec[T any](f func() (T, error)) Future[T] {
	var (
		result T
		err    error
	)
	c := make(chan struct{})
	go func() {
		defer close(c)

		result, err = f()
	}()
	return future[T]{
		await: func(ctx context.Context) (T, error) {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-c:
				return result, err
			}
		},
	}
}
