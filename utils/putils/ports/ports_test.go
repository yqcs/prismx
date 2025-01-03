package ports

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	t.Run("valid-ports-strings", func(t *testing.T) {
		ports := []interface{}{"1", "10000", "65535", 1, 10000, 65535}
		for _, port := range ports {
			require.True(t, IsValid(port))
		}
	})
	t.Run("invalid-ports", func(t *testing.T) {
		ports := []interface{}{"", "-1", "0", "65536", 0, -1, 65536, 2.1, "a"}
		for _, port := range ports {
			require.False(t, IsValid(port))
		}
	})
}
