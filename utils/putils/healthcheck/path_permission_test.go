package healthcheck

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPathPermission(t *testing.T) {
	t.Run("file with read and write permissions", func(t *testing.T) {
		filename := "testfile_read_write.txt"
		_, err := os.Create(filename)
		defer os.Remove(filename)
		assert.NoError(t, err)

		permission := CheckPathPermission(filename)
		assert.NoError(t, permission.Error)
		assert.Equal(t, true, permission.isReadable)
		assert.Equal(t, true, permission.isWritable)
	})

	t.Run("non-existing file", func(t *testing.T) {
		filename := "non_existing_file.txt"
		permission := CheckPathPermission(filename)
		assert.Error(t, permission.Error)
	})

	t.Run("file without write permission", func(t *testing.T) {
		filename := "testfile_read_only.txt"
		file, err := os.Create(filename)
		assert.NoError(t, err)

		err = file.Chmod(0444) // read-only permissions
		assert.NoError(t, err)

		defer os.Remove(filename)
		permission := CheckPathPermission(filename)

		assert.NoError(t, permission.Error)
		assert.Equal(t, true, permission.isReadable)
		assert.Equal(t, false, permission.isWritable)
	})
}
