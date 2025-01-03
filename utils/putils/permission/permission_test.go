//go:build windows || linux

package permissionutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsRoot(t *testing.T) {
	isRoot, err := checkCurrentUserRoot()
	require.Nil(t, err)
	require.NotNil(t, isRoot)
}
