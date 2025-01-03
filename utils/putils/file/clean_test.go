package fileutil

import (
	"io"
	"os"
	"strings"
	"testing"
)

func FuzzSafeOpen(f *testing.F) {

	// ==========setup==========

	bin, err := os.ReadFile("tests/path-traversal.txt")
	if err != nil {
		f.Fatalf("failed to read file: %s", err)
	}

	fuzzPayloads := strings.Split(string(bin), "\n")

	file, err := os.CreateTemp("", "*")
	if err != nil {
		f.Fatal(err)
	}
	_, _ = file.WriteString("pwned!")
	_ = file.Close()

	defer func(tmp string) {
		if err = os.Remove(tmp); err != nil {
			panic(err)
		}
	}(file.Name())

	// ==========fuzzing==========

	for _, payload := range fuzzPayloads {
		f.Add(strings.ReplaceAll(payload, "{FILE}", f.Name()), f.Name())

	}
	f.Fuzz(func(t *testing.T, fuzzPath string, targetPath string) {
		cleaned, err := CleanPath(fuzzPath)
		if err != nil {
			// Ignore errors
			return
		}
		if cleaned != targetPath {
			// cleaned path is different from target file
			// so verify if 'path' is actually valid and not random chars
			result, err := SafeOpen(cleaned)
			if err != nil {
				// Ignore errors
				return
			}
			defer result.Close()
			bin, _ := io.ReadAll(result)
			if string(bin) == "pwned!" {
				t.Fatalf("pwned! cleaned=%s ,input=%s", cleaned, fuzzPath)
			}
		}

	})
}
