//go:build linux || darwin

package permissionutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// run the file permission tests on linux and osx
func TestFilePermissions(t *testing.T) {
	t.Run("TestFileAllReadWriteExecute", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		require.Nil(t, err, "Failed to create test file: %v", err)
		testFileName := file.Name()
		defer os.Remove(testFileName)
		defer file.Close()

		// Set the file permissions
		err = file.Chmod(os.FileMode(AllReadWriteExecute))
		require.Nil(t, err, "Failed to set file permissions: %v", err)

		// Get the file permissions
		fileInfo, err := os.Stat(testFileName)
		require.Nil(t, err, "Failed to get file info: %v", err)
		// Check if the file permissions match the defined constants
		match := fileInfo.Mode().Perm().String() != "-rwxrwxrwx" || fileInfo.Mode().Perm() != os.FileMode(AllReadWriteExecute)
		require.False(t, match, "File permissions do not match. Expected: %s, Actual: %s", os.FileMode(AllReadWriteExecute).String(), fileInfo.Mode().Perm().String())
	})

	t.Run("TestFileUserReadWriteExecute", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		require.Nil(t, err, "Failed to create test file: %v", err)

		testFileName := file.Name()
		defer os.Remove(testFileName)
		defer file.Close()

		// Set the file permissions
		err = file.Chmod(os.FileMode(UserReadWriteExecute))
		require.Nil(t, err, "Failed to set file permissions: %v", err)

		// Get the file permissions
		fileInfo, err := os.Stat(testFileName)
		require.Nil(t, err, "Failed to get file info: %v", err)

		// Check if the file permissions match the defined constants
		match := fileInfo.Mode().Perm().String() != "-rwx------" || fileInfo.Mode().Perm() != os.FileMode(UserReadWriteExecute)
		require.False(t, match, "File permissions do not match. Expected: %s, Actual: %s", os.FileMode(UserReadWriteExecute).String(), fileInfo.Mode().Perm().String())
	})

	t.Run("TestFileGroupReadWriteExecute", func(t *testing.T) {
		file, err := os.CreateTemp("", "testfile")
		require.Nil(t, err, "Failed to create test file: %v", err)

		testFileName := file.Name()
		defer os.Remove(testFileName)
		defer file.Close()

		// Set the file permissions
		err = file.Chmod(os.FileMode(UserReadWriteExecute | GroupReadWriteExecute))
		require.Nil(t, err, "Failed to set file permissions: %v", err)

		// Get the file permissions
		fileInfo, err := os.Stat(testFileName)
		require.Nil(t, err, "Failed to get file info: %v", err)

		// Check if the file permissions match the defined constants
		match := fileInfo.Mode().Perm().String() != "-rwxrwx---" || fileInfo.Mode().Perm() != os.FileMode(UserReadWriteExecute|GroupReadWriteExecute)
		require.False(t, match, "File permissions do not match. Expected: %s, Actual: %s", os.FileMode(UserReadWriteExecute|GroupReadWriteExecute).String(), fileInfo.Mode().Perm().String())
	})
}

func TestUpdateFilePerm(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	require.Nil(t, err, "Error creating temp file: %v", err)
	defer os.Remove(tempFile.Name())

	// Set the desired file permissions
	expectedPermissions := AllReadWrite

	err = UpdateFilePerm(tempFile.Name(), expectedPermissions)
	require.Nil(t, err, "Error updating file permissions: %v", err)

	// Get the updated file information
	updatedFileInfo, err := os.Stat(tempFile.Name())
	require.Nil(t, err, "Error getting updated file information: %v", err)

	// Check if the updated file permissions match expected permissions
	updatedFileMode := updatedFileInfo.Mode().Perm()
	match := updatedFileMode != os.FileMode(expectedPermissions)
	require.False(t, match, "Invalid file permissions, expected: %v, got: %v", os.FileMode(expectedPermissions), updatedFileMode)
}
