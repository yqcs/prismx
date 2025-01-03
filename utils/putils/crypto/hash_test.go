package cryptoutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSHA256Sum(t *testing.T) {
	tests := map[string]string{
		"test":  "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		"test1": "1b4f0e9851971998e732078544c96b36c3d01cedf7caa332359d6f1d83567014",
	}
	for item, hash := range tests {
		require.Equal(t, hash, SHA256Sum(item), "hash is different")
	}
}
