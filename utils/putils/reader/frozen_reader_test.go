package reader

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestFrozenReader(t *testing.T) {
	forever := func() {
		wrappedStdin := FrozenReader{}
		_, err := io.Copy(os.Stdout, wrappedStdin)
		if err != nil {
			return
		}
	}
	go forever()
	<-time.After(10 * time.Second)
}
