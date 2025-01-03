package healthcheck

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectEnvironmentInfo(t *testing.T) {
	t.Run("Collect Environment Info", func(t *testing.T) {
		programVersion := "1.0.0"

		environmentInfo := CollectEnvironmentInfo(programVersion)
		assert.NoError(t, environmentInfo.Error, "Error should not have occurred when collecting environment info")
		assert.NotNil(t, environmentInfo, "EnvironmentInfo should not be nil")
		assert.Equal(t, programVersion, environmentInfo.ProgramVersion, "Program version should match program version")
		assert.Equal(t, runtime.GOARCH, environmentInfo.Arch, "Architecture should match runtime")
		assert.Equal(t, runtime.Compiler, environmentInfo.Compiler, "Compiler should match runtime")
		assert.Equal(t, runtime.Version(), environmentInfo.GoVersion, "Go version should match runtime")
		assert.Equal(t, runtime.GOOS, environmentInfo.OSName, "OS name should match runtime")
		assert.Equal(t, os.Getenv("PATH"), environmentInfo.PathEnvVar, "PATH environment variable should match system PATH")
	})
}
