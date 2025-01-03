package reader

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeoutReader(t *testing.T) {
	wrappedStdin := TimeoutReader{
		Reader:  FrozenReader{},
		Timeout: time.Duration(2 * time.Second),
	}
	_, err := io.Copy(os.Stdout, wrappedStdin)
	require.NotNil(t, err)
}
