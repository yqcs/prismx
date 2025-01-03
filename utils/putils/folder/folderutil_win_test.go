//go:build windows

package folderutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathInfo(t *testing.T) {
	got, err := NewPathInfo("c:\\a\\b\\c")
	assert.Nil(t, err)
	gotPaths, err := got.Paths()
	assert.Nil(t, err)
	assert.EqualValues(t, []string{".", "c:\\", "c:\\a", "c:\\a\\b", "c:\\a\\b\\c"}, gotPaths)
	gotMeshPaths, err := got.MeshWith("test.txt")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"test.txt", "c:\\test.txt", "c:\\a\\test.txt", "c:\\a\\b\\test.txt", "c:\\a\\b\\c\\test.txt"}, gotMeshPaths)
}
