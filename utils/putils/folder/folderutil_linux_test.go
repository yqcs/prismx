//go:build !windows

package folderutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathInfo(t *testing.T) {
	got, err := NewPathInfo("/a/b/c")
	assert.Nil(t, err)
	gotPaths, err := got.Paths()
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"/", "/a", "/a/b", "/a/b/c"}, gotPaths)
	gotMeshPaths, err := got.MeshWith("test.txt")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"/test.txt", "/a/test.txt", "/a/b/test.txt", "/a/b/c/test.txt"}, gotMeshPaths)
}
