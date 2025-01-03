package folderutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fileutil "prismx_cli/utils/putils/file"
	osutils "prismx_cli/utils/putils/os"
)

func TestGetFiles(t *testing.T) {
	// get files from current folder
	files, err := GetFiles(".")
	require.Nilf(t, err, "couldn't retrieve the list of files: %s", err)

	// we check only if the number of files is bigger than zero
	require.Positive(t, len(files), "no files could be retrieved: %s", err)
}

func TestSyncDirectory(t *testing.T) {
	t.Run("destination folder creation error", func(t *testing.T) {
		err := SyncDirectory("/source", "/:/dest")
		assert.Error(t, err)
	})

	t.Run("source folder not found error", func(t *testing.T) {
		err := SyncDirectory("/notExistingFolder", "/dest")
		assert.Error(t, err)
	})

	t.Run("source and destination are the same", func(t *testing.T) {
		// setup
		// some files in a temp dir
		sourceDir := t.TempDir()
		defer os.RemoveAll(sourceDir)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file1.txt"), []byte("file1"), os.ModePerm)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file2.txt"), []byte("file2"), os.ModePerm)

		// when: try to migrate files
		err := SyncDirectory(sourceDir, sourceDir)

		// then: verify if files migrated successfully
		assert.Error(t, err)

		assert.True(t, fileutil.FileExists(filepath.Join(sourceDir, "/file1.txt")))
		assert.True(t, fileutil.FileExists(filepath.Join(sourceDir, "/file2.txt")))
	})

	t.Run("successful migration with source dir removal", func(t *testing.T) {
		// setup
		// some files in a temp dir
		sourceDir := t.TempDir()
		defer os.RemoveAll(sourceDir)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file1.txt"), []byte("file1"), os.ModePerm)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file2.txt"), []byte("file2"), os.ModePerm)
		_ = os.Mkdir(filepath.Join(sourceDir, "/dir1"), os.ModePerm)
		_ = os.WriteFile(filepath.Join(sourceDir, "/dir1", "/file3.txt"), []byte("file3"), os.ModePerm)
		_ = os.Mkdir(filepath.Join(sourceDir, "/dir2"), os.ModePerm)

		// destination directory
		destinationDir := t.TempDir()
		defer os.RemoveAll(destinationDir)

		// when: try to migrate files
		err := SyncDirectory(sourceDir, destinationDir)

		// then: verify if files migrated successfully
		assert.NoError(t, err, sourceDir, destinationDir)

		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/file1.txt")))
		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/file2.txt")))
		assert.True(t, fileutil.FolderExists(filepath.Join(destinationDir, "/dir1")))
		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/dir1", "/file3.txt")))

		assert.False(t, fileutil.FolderExists(filepath.Join(destinationDir, "/dir2")))

		assert.False(t, fileutil.FolderExists(sourceDir))
	})

	t.Run("successful migration without source dir removal", func(t *testing.T) {
		// setup
		// some files in a temp dir
		sourceDir := t.TempDir()
		defer os.RemoveAll(sourceDir)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file1.txt"), []byte("file1"), os.ModePerm)
		_ = os.WriteFile(filepath.Join(sourceDir, "/file2.txt"), []byte("file2"), os.ModePerm)
		_ = os.Mkdir(filepath.Join(sourceDir, "/dir1"), os.ModePerm)
		_ = os.WriteFile(filepath.Join(sourceDir, "/dir1", "/file3.txt"), []byte("file3"), os.ModePerm)
		_ = os.Mkdir(filepath.Join(sourceDir, "/dir2"), os.ModePerm)

		// destination directory
		destinationDir := t.TempDir()
		defer os.RemoveAll(destinationDir)

		// when: try to migrate files
		RemoveSourceDirAfterSync = false
		err := SyncDirectory(sourceDir, destinationDir)

		// then: verify if files migrated successfully
		assert.NoError(t, err)

		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/file1.txt")))
		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/file2.txt")))
		assert.True(t, fileutil.FolderExists(filepath.Join(destinationDir, "/dir1")))
		assert.True(t, fileutil.FileExists(filepath.Join(destinationDir, "/dir1", "/file3.txt")))

		assert.False(t, fileutil.FolderExists(filepath.Join(destinationDir, "/dir2")))

		assert.True(t, fileutil.FolderExists(sourceDir))
	})
}

func TestIsWritable(t *testing.T) {
	t.Run("Test writable directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test-dir")
		assert.NoError(t, err)
		defer os.RemoveAll(tempDir)

		assert.True(t, IsWritable(tempDir), "expected directory to be writable")
	})

	t.Run("Test non-existent directory", func(t *testing.T) {
		nonExistentDir := "/path/to/non/existent/dir"
		assert.False(t, IsWritable(nonExistentDir), "expected directory to not be writable")
	})

	t.Run("Test non-writable directory", func(t *testing.T) {
		// on windows bitsets are applied only to files
		// https://github.com/golang/go/issues/35042
		if osutils.IsWindows() {
			return
		}

		nonWritableDir := "non-writable-dir"
		err := os.Mkdir(nonWritableDir, 0400)
		assert.NoError(t, err)
		defer os.RemoveAll(nonWritableDir)

		// Make the directory non-writable.
		err = os.Chmod(nonWritableDir, 0400)
		assert.NoError(t, err)

		assert.False(t, IsWritable(nonWritableDir), "expected directory to not be writable")
	})

	t.Run("Test with a file instead of a directory", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "test-file")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		assert.False(t, IsWritable(tempFile.Name()), "expected file to not be considered a writable directory")
	})
}
