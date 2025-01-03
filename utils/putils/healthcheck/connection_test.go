package healthcheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckConnection(t *testing.T) {
	t.Run("Test successful connection", func(t *testing.T) {
		info := CheckConnection("scanme.sh", 80, "tcp", 1*time.Second)
		assert.NoError(t, info.Error)
		assert.True(t, info.Successful)
		assert.Equal(t, "scanme.sh", info.Host)
		assert.Contains(t, info.Message, "Successful")
	})

	t.Run("Test unsuccessful connection", func(t *testing.T) {
		info := CheckConnection("invalid.website", 80, "tcp", 1*time.Second)
		assert.Error(t, info.Error)
	})

	t.Run("Test timeout connection", func(t *testing.T) {
		info := CheckConnection("192.0.2.0", 80, "tcp", 1*time.Millisecond)
		assert.Error(t, info.Error)
	})
}
