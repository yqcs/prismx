package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAsync(t *testing.T) {
	// Async
	do := Exec(func() (bool, error) {
		time.Sleep(2 * time.Second)
		return true, nil
	})

	// do some other stuff
	time.Sleep(time.Second)

	// Await
	ok, err := do.Await()
	require.Nil(t, err)
	require.True(t, ok)
}
