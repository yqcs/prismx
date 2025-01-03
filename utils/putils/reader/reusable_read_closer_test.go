package reader

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReusableReader(t *testing.T) {
	testcases := []interface{}{
		strings.NewReader("test"),
		bytes.NewBuffer([]byte("test")),
		bytes.NewBufferString("test"),
		bytes.NewReader([]byte("test")),
		[]byte("test"),
		"test",
	}
	for _, v := range testcases {
		t.Run("sequential reuse", func(t *testing.T) {
			reusableReader, err := NewReusableReadCloser(v)
			require.Nil(t, err)

			for i := 0; i < 100; i++ {
				n, err := io.Copy(io.Discard, reusableReader)
				require.Nil(t, err)
				require.Positive(t, n)

				bin, err := io.ReadAll(reusableReader)
				require.Nil(t, err)
				require.Len(t, bin, 4)
			}
		})

		// todo: readers shouldn't be used concurrently, so here we just try to catch pontential read concurring with resets panics
		t.Run("concurrent-reset-with-read", func(t *testing.T) {
			reusableReader, err := NewReusableReadCloser(v)
			require.Nil(t, err)

			var wg sync.WaitGroup

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					_, err := io.Copy(io.Discard, reusableReader)
					require.Nil(t, err)
					_, err = io.ReadAll(reusableReader)
					require.Nil(t, err)
				}()
			}

			wg.Wait()
		})
	}
}
