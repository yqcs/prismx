package generic

import (
	"sync"
	"testing"
)

func TestDo(t *testing.T) {
	val := 10
	l := WithLock(val)
	l.Do(func(v int) {
		if v != val {
			t.Errorf("Expected %d, got %d", val, v)
		}
	})
}

func TestLockableConcurrency(t *testing.T) {
	l := WithLock(0)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				l.Do(func(v int) {
					v++
					l.V = v
				})
			}
		}()
	}

	wg.Wait()

	if l.V != 100*1000 {
		t.Errorf("Expected counter to be %d, but got %d", 100*1000, l.V)
	}
}

func TestLockableStringManipulation(t *testing.T) {
	str := "initial"
	l := WithLock(str)

	l.Do(func(s string) {
		s += " - updated"
		l.V = s
	})

	if l.V != "initial - updated" {
		t.Errorf("Expected 'initial - updated', got '%s'", str)
	}
}
