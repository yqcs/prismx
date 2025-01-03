package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBumpVersionCLI(t *testing.T) {
	testCases := []struct {
		name           string
		fileContent    string
		varName        string
		part           string
		expectedErr    bool
		expectedOutput string
	}{
		{
			name:           "Test v1.0.0",
			fileContent:    "package main\n\nvar version = \"v1.0.0\"",
			varName:        "version",
			part:           "patch",
			expectedErr:    false,
			expectedOutput: "Bump from v1.0.0 to v1.0.1\n",
		},
		{
			name:           "Test 1.0.8",
			fileContent:    "package main\n\nvar version = \"1.0.8\"",
			varName:        "version",
			part:           "patch",
			expectedErr:    false,
			expectedOutput: "Bump from 1.0.8 to v1.0.9\n",
		},
		{
			name:           "Test v1.0.9-dev",
			fileContent:    "package main\n\nvar version = \"v1.0.9-dev\"",
			varName:        "version",
			part:           "patch",
			expectedErr:    false,
			expectedOutput: "Bump from v1.0.9-dev to v1.0.9\n",
		},
		{
			name:        "Test 1.0.85.1",
			fileContent: "package main\n\nvar version = \"1.0.85.1\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unwanted prefixes/suffixes: with x prefix",
			fileContent: "package main\n\nvar version = \"x1.0.0\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unwanted prefixes/suffixes: with x suffix",
			fileContent: "package main\n\nvar version = \"1.0.0x\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unwanted prefixes/suffixes: with new lines",
			fileContent: "package main\n\nvar version = \"\n\n1.0.0\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unwanted prefixes/suffixes: with space",
			fileContent: "package main\n\nvar version = \"  	1.0.0\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unwanted prefixes/suffixes: with multiple dot",
			fileContent: "package main\n\nvar version = \"1.0..0\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with negative numbers",
			fileContent: "package main\n\nvar version = \"-1.0.0\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with unparseable numbers",
			fileContent: "package main\n\nvar version = \"1.0.X\"",
			varName:     "version",
			part:        "patch",
			expectedErr: true,
		},
		{
			name:        "Test with big numbers",
			fileContent: "package main\n\nvar version = \"1.0.111111111111111111111111111111111111111111111111111111111111111111111111\"",
			varName:     "minor",
			part:        "patch",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file and write the content to it
			tempFile, err := os.CreateTemp(os.TempDir(), "prefix-")
			if err != nil {
				t.Fatalf("Cannot create temporary file: %s", err)
			}

			defer os.Remove(tempFile.Name())

			_, err = tempFile.Write([]byte(tc.fileContent))
			if err != nil {
				t.Fatalf("Failed to write to temporary file: %s", err)
			}

			cmd := exec.Command("go", "run", "versionbump.go", "-file", tempFile.Name(), "-var", tc.varName, "-part", tc.part)

			output, err := cmd.CombinedOutput()
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Contains(t, string(output), tc.expectedOutput)
			}
		})
	}
}
