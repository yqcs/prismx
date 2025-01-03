package syscallutil

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	osutils "prismx_cli/utils/putils/os"
)

func TestLoadLibrary(t *testing.T) {
	t.Run("Test valid library", func(t *testing.T) {
		var lib string
		if osutils.IsWindows() {
			lib = "ucrtbase.dll"
		} else if osutils.IsOSX() {
			lib = "libSystem.dylib"
		} else if osutils.IsLinux() {
			lib = "libc.so.6"
		} else {
			panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
		}

		_, err := LoadLibrary(lib)
		require.NoError(t, err, "should not return an error for valid library")
	})

	t.Run("Test invalid library", func(t *testing.T) {
		var lib string
		if osutils.IsWindows() {
			lib = "C:\\path\\to\\invalid\\library.dll"
		} else if osutils.IsOSX() {
			lib = "/path/to/invalid/library.dylib"
		} else if osutils.IsLinux() {
			lib = "/path/to/invalid/library.so"
		} else {
			panic(fmt.Errorf("GOOS=%s is not supported", runtime.GOOS))
		}

		_, err := LoadLibrary(lib)
		require.Error(t, err, "should return an error for invalid library")
	})
}
