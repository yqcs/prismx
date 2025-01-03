package executil

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var newLineMarker string

func init() {
	if runtime.GOOS == "windows" {
		newLineMarker = "\r\n"
	} else {
		newLineMarker = "\n"
	}
}

func TestRun(t *testing.T) {
	// try to run the echo command
	s, err := Run("echo test")
	require.Nil(t, err, "failed execution", err)
	require.Equal(t, "test"+newLineMarker, s, "output doesn't contain expected result", s)
}

func TestRunAdv(t *testing.T) {
	testcases := []struct {
		GOOS     string // OS
		Command  string
		Expected string // expected output
		Contains string // expected output contains
	}{
		// Tests With Flags
		{"darwin", "uname -s", "Darwin", ""},
		{"linux", "uname -s", "Linux", ""},
		{"windows", "cmd /c ver", "", "Windows"},
		// Tests With CMD PIPE
		{"windows", `systeminfo | findstr /B /C:"OS Name"`, "", "Windows"},
		{"darwin", `sw_vers | grep -i "ProductName"`, "", "macOS"},
		{"linux", `uname -a  | cut -d " " -f 1`, "Linux", ""},
		// Other Shell Specific Features
		{"windows", `cmd /c " echo This && echo Works"`, "This \r\nWorks", ""},
		{"linux", "true && echo This Works", "This Works", ""},
		{"darwin", "true && echo This Works", "This Works", ""},
	}

	runFunc := func(cmd string, expected string, contains string) {
		s, err := Run(cmd)
		require.Nilf(t, err, "%v failed to execute", cmd)
		if expected != "" {
			require.Equal(t, expected+newLineMarker, s)
		} else if contains != "" {
			require.Contains(t, s, contains)
		} else {
			t.Logf("Malformed test case : %v", cmd)
		}
		t.Logf("Test Successful: %v", cmd)
	}

	for _, v := range testcases {
		switch v.GOOS {
		case "windows":
			if runtime.GOOS != "windows" {
				continue
			}
			runFunc(v.Command, v.Expected, v.Contains)
		case "darwin":
			if runtime.GOOS != "darwin" {
				continue
			}
			runFunc(v.Command, v.Expected, v.Contains)
		case "linux":
			if runtime.GOOS != "linux" {
				continue
			}
			runFunc(v.Command, v.Expected, v.Contains)
		default:
			t.Logf("No Unit Test Available for this platform")

		}
	}
}

func TestRunSafe(t *testing.T) {
	_, err := RunSafe(`whoami | grep Hello`)
	require.Error(t, err)
}

func TestRunSh(t *testing.T) {
	if runtime.GOOS == "windows" {
		return
	}
	// try to run the echo command
	s, err := RunSh("echo", "test")
	require.Nil(t, err, "failed execution", err)
	require.Equal(t, "test"+newLineMarker, s, "output doesn't contain expected result", s)
}

func TestRunPS(t *testing.T) {
	if runtime.GOOS != "windows" {
		return
	}
	// run powershell command (runs in both ps1 and ps2)
	s, err := RunPS("get-host")
	require.Nil(t, err, "failed execution", err)
	require.Contains(t, s, "Microsoft.PowerShell", "failed to run powershell command get-host")
}
