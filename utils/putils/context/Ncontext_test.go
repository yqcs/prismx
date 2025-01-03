package contextutil_test

import (
	"context"
	"errors"
	"testing"
	"time"

	contextutil "prismx_cli/utils/putils/context"
)

func TestExecFuncWithTwoReturns(t *testing.T) {
	t.Run("function completes before context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		fn := func() (int, error) {
			time.Sleep(1 * time.Second)
			return 42, nil
		}

		val, err := contextutil.ExecFuncWithTwoReturns(ctx, fn)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if val != 42 {
			t.Errorf("Unexpected return value: got %v, want 42", val)
		}
	})

	t.Run("context cancelled before function completes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		fn := func() (int, error) {
			time.Sleep(2 * time.Second)
			return 42, nil
		}

		_, err := contextutil.ExecFuncWithTwoReturns(ctx, fn)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context deadline exceeded error, got: %v", err)
		}
	})
}

func TestExecFuncWithThreeReturns(t *testing.T) {
	t.Run("function completes before context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		fn := func() (int, string, error) {
			time.Sleep(1 * time.Second)
			return 42, "hello", nil
		}

		val1, val2, err := contextutil.ExecFuncWithThreeReturns(ctx, fn)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if val1 != 42 || val2 != "hello" {
			t.Errorf("Unexpected return values: got %v and %v, want 42 and 'hello'", val1, val2)
		}
	})

	t.Run("context cancelled before function completes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		fn := func() (int, string, error) {
			time.Sleep(2 * time.Second)
			return 42, "hello", nil
		}

		_, _, err := contextutil.ExecFuncWithThreeReturns(ctx, fn)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Expected context deadline exceeded error, got: %v", err)
		}
	})
}

func TestExecFunc(t *testing.T) {
	t.Run("function completes before context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		fn := func() {
			time.Sleep(1 * time.Second)
		}

		err := contextutil.ExecFunc(ctx, fn)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("context cancelled before function completes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		fn := func() {
			time.Sleep(2 * time.Second)
		}

		err := contextutil.ExecFunc(ctx, fn)
		if err != context.DeadlineExceeded {
			t.Errorf("Expected context deadline exceeded error, got: %v", err)
		}
	})
}
