package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFileOrFolderExists(t *testing.T) {
	tests := map[string]bool{
		"file.go": true,
		"aaa.bbb": false,
		".":       true,
		"../file": true,
		"aabb":    false,
	}
	for fpath, mustExist := range tests {
		exist := FileOrFolderExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\": %v", fpath, exist)
	}
}

func TestFileExists(t *testing.T) {
	tests := map[string]bool{
		"file.go": true,
		"aaa.bbb": false,
		"/":       false,
	}
	for fpath, mustExist := range tests {
		exist := FileExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\": %v", fpath, exist)
	}
}

func TestFolderExists(t *testing.T) {
	tests := map[string]bool{
		".":       true,
		"../file": true,
		"aabb":    false,
	}
	for fpath, mustExist := range tests {
		exist := FolderExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\"", fpath)
	}
}

func TestDeleteFilesOlderThan(t *testing.T) {
	// create a temporary folder with a couple of files
	fo, err := os.MkdirTemp("", "")
	require.Nil(t, err, "couldn't create folder: %s", err)
	ttl := time.Duration(1 * time.Second)
	sleepTime := time.Duration(3 * time.Second)

	// defer temporary folder removal
	defer os.RemoveAll(fo)
	checkFolderErr := func(err error) {
		require.Nil(t, err, "couldn't create folder: %s", err)
	}
	checkFiles := func(fileName string) {
		require.False(t, FileExists(fileName), "file \"%s\" still exists", fileName)
	}
	createFile := func() string {
		fi, err := os.CreateTemp(fo, "")
		require.Nil(t, err, "couldn't create f: %s", err)
		fName := fi.Name()
		fi.Close()
		return fName
	}
	t.Run("prefix props test", func(t *testing.T) {
		fName := createFile()
		fileInfo, _ := os.Stat(fName)
		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan: ttl,
			Prefix:    fileInfo.Name(),
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		checkFiles(fName)
	})
	t.Run("suffix props test", func(t *testing.T) {
		fName := createFile()
		fileInfo, _ := os.Stat(fName)
		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan: ttl,
			Suffix:    string(fileInfo.Name()[len(fileInfo.Name())-1]),
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		checkFiles(fName)
	})
	t.Run("regex pattern props test", func(t *testing.T) {
		fName := createFile()
		fName1 := createFile()

		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan:    ttl,
			RegexPattern: "[0-9]",
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		checkFiles(fName)
		checkFiles(fName1)
	})
	t.Run("custom check props test", func(t *testing.T) {
		fName := createFile()
		fName1 := createFile()

		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan: ttl,
			CustomCheck: func(filename string) bool {
				return true
			},
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		checkFiles(fName)
		checkFiles(fName1)
	})
	t.Run("custom check props negative test", func(t *testing.T) {
		fName := createFile()
		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan: ttl,
			CustomCheck: func(filename string) bool {
				return false // should not delete the file
			},
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		require.True(t, FileExists(fName), "file \"%s\" should exists", fName)
	})
	t.Run("callback props test", func(t *testing.T) {
		fName := createFile()
		fName1 := createFile()

		// sleep for 5 seconds
		time.Sleep(sleepTime)
		// delete files older than 5 seconds
		filter := FileFilters{
			OlderThan: ttl,
			CustomCheck: func(filename string) bool {
				return true
			},
			Callback: func(filename string) error {
				t.Log("deleting file manually")
				return os.Remove(filename)
			},
		}
		err = DeleteFilesOlderThan(fo, filter)
		checkFolderErr(err)
		checkFiles(fName)
		checkFiles(fName1)
	})
}

func TestDownloadFile(t *testing.T) {
	// attempt to download http://ipv4.download.thinkbroadband.com/5MB.zip to temp folder
	tmpfile, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create folder: %s", err)
	fname := tmpfile.Name()

	os.Remove(fname)

	err = DownloadFile(fname, "http://ipv4.download.thinkbroadband.com/5MB.zip")
	require.Nil(t, err, "couldn't download file: %s", err)

	require.True(t, FileExists(fname), "file \"%s\" doesn't exists", fname)

	// remove the downloaded file
	os.Remove(fname)
}

func tmpFolderName(s string) string {
	return filepath.Join(os.TempDir(), s)
}

func TestCreateFolders(t *testing.T) {
	tests := []string{
		tmpFolderName("a"),
		tmpFolderName("b"),
	}
	err := CreateFolders(tests...)
	require.Nil(t, err, "couldn't download file: %s", err)

	for _, folder := range tests {
		fexists := FolderExists(folder)
		require.True(t, fexists, "folder %s doesn't exist", fexists)
	}

	// remove folders
	for _, folder := range tests {
		os.Remove(folder)
	}
}

func TestCreateFolder(t *testing.T) {
	tst := tmpFolderName("a")
	err := CreateFolder(tst)
	require.Nil(t, err, "couldn't download file: %s", err)

	fexists := FolderExists(tst)
	require.True(t, fexists, "folder %s doesn't exist", fexists)

	os.Remove(tst)
}

func TestHasStdin(t *testing.T) {
	require.False(t, HasStdin(), "stdin in test")
}

func TestReadFile(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)

	fileContentLines := strings.Split(fileContent, "\n")
	// compare file lines
	c, err := ReadFile(fname)
	require.Nil(t, err, "couldn't open file: %s", err)
	i := 0
	for line := range c {
		require.Equal(t, fileContentLines[i], line, "lines don't match")
		i++
	}
}

func TestReadFileWithBufferSize(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)

	fileContentLines := strings.Split(fileContent, "\n")
	// compare file lines
	c, err := ReadFileWithBufferSize(fname, 64*1024)
	require.Nil(t, err, "couldn't open file: %s", err)
	i := 0
	for line := range c {
		require.Equal(t, fileContentLines[i], line, "lines don't match")
		i++
	}
}

func TestPermissions(t *testing.T) {
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	f.Close()
	defer os.Remove(fname)

	ok, err := IsReadable(fname)
	require.True(t, ok)
	require.Nil(t, err)
	ok, err = IsWriteable(fname)
	require.True(t, ok)
	require.Nil(t, err)
}

func TestUseMusl(t *testing.T) {
	executablePath, err := os.Executable()
	require.Nil(t, err)
	_, err = UseMusl(executablePath)
	switch runtime.GOOS {
	case "windows", "darwin":
		require.NotNil(t, err)
	case "linux":
		require.Nil(t, err)
	}
}

func TestReadFileWithReader(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "output")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)
	fileContentLines := strings.Split(fileContent, "\n")
	f, err = os.Open(fname)
	require.Nil(t, err, "couldn't create file: %s", err)
	// compare file lines
	c, _ := ReadFileWithReader(f)
	i := 0
	for line := range c {
		require.Equal(t, fileContentLines[i], line, "lines don't match")
		i++
	}
	f.Close()
}

func TestReadFileWithReaderAndBufferSize(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)
	fileContentLines := strings.Split(fileContent, "\n")
	f, err = os.Open(fname)
	require.Nil(t, err, "couldn't create file: %s", err)
	// compare file lines
	c, _ := ReadFileWithReaderAndBufferSize(f, 64*1024)
	i := 0
	for line := range c {
		require.Equal(t, fileContentLines[i], line, "lines don't match")
		i++
	}
	f.Close()
}

func TestCopyFile(t *testing.T) {
	fileContent := `test
	test1
	test2`
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(fileContent))
	f.Close()
	defer os.Remove(fname)
	fnameCopy := fmt.Sprintf("%s-copy", f.Name())
	err = CopyFile(fname, fnameCopy)
	require.Nil(t, err, "couldn't copy file: %s", err)
	require.True(t, FileExists(fnameCopy), "file \"%s\" doesn't exists", fnameCopy)
	os.Remove(fnameCopy)
}

func TestGetTempFileName(t *testing.T) {
	fname, _ := GetTempFileName()
	defer os.Remove(fname)
	require.NotEmpty(t, fname)
}

func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	var ts TestStruct
	err := Unmarshal(JSON, []byte(`{"name":"test"}`), &ts)
	require.Nil(t, err)
	require.Equal(t, "test", ts.Name)
}

func TestMarshal(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	ts := TestStruct{Name: "test"}
	fs, err := GetTempFileName()
	require.Nil(t, err)
	defer RemoveAll(fs)
	err = Marshal(JSON, []byte(fs), ts)
	require.Nil(t, err)
	data, err := os.ReadFile(fs)
	require.Nil(t, err)
	require.Equal(t, `{"name":"test"}`, strings.TrimSpace(string(data)))
}

func TestUnmarshalFromReader(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	var ts TestStruct
	err := UnmarshalFromReader(JSON, strings.NewReader(`{"name":"test"}`), &ts)
	require.Nil(t, err)
	require.Equal(t, "test", ts.Name)
}

func TestMarshalToWriter(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	ts := TestStruct{Name: "test"}
	var data []byte
	buf := bytes.NewBuffer(data)
	err := MarshalToWriter(JSON, buf, ts)
	require.Nil(t, err)
	require.Equal(t, `{"name":"test"}`, strings.TrimSpace(buf.String()))
}

func TestExecutableName(t *testing.T) {
	require.NotEmpty(t, ExecutableName())
}

func TestRemoveAll(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "")
	require.Nil(t, err, "couldn't create folder: %s", err)
	f, err := os.CreateTemp(tmpdir, "")
	require.Nil(t, err, "couldn't create file: %s", err)
	f.Close()
	errs := RemoveAll(tmpdir)
	require.Equal(t, 0, len(errs), "couldn't remove folder: %s", errs)
}

func TestCountLineWithSeparator(t *testing.T) {
	testcases := []struct {
		filename       string
		expectedLines  uint
		shouldError    bool
		expectedError  string
		skipEmptyLines bool
		separator      string
	}{
		{
			filename:      "tests/standard.txt",
			expectedLines: 5,
			separator:     "\n",
		},
		{
			filename:      "tests/empty_lines.txt",
			expectedLines: 18,
			separator:     "\n",
		},
		{
			filename:      "tests/pipe_separator.txt",
			expectedLines: 5,
			separator:     "|",
		},
		{
			filename:      "nonexistent.txt",
			shouldError:   true,
			expectedLines: 0,
			separator:     "\n",
		},
		{
			filename:      "tests/standard.txt",
			separator:     "",
			shouldError:   true,
			expectedError: "invalid separator",
		},
	}
	for _, test := range testcases {
		linesCount, err := CountLinesWithSeparator([]byte(test.separator), test.filename)
		if test.shouldError {
			require.NotNil(t, err)
			if test.expectedError != "" {
				require.EqualError(t, err, test.expectedError)
			}
		} else {
			require.Nil(t, err)
			require.Equal(t, test.expectedLines, linesCount)
		}
	}
}

func TestSubstituteConfigFromEnvVars(t *testing.T) {
	configFileContent := `test:
	- id: some_id
	  channel: $CHANNEL
	  username: $USER
	  webhook_url: $WEBHOOK
	  threads: $THREADS
  `
	f, err := os.CreateTemp("", "")
	require.Nil(t, err, "couldn't create file: %s", err)
	fname := f.Name()
	_, _ = f.Write([]byte(configFileContent))
	f.Close()
	defer os.Remove(fname)

	os.Setenv("CHANNEL", "test_channel")
	os.Setenv("USER", "test_user")
	os.Setenv("WEBHOOK", "test_webhook")
	os.Setenv("THREADS", "test_threads")

	expectedFileContent := `test:
	- id: some_id
	  channel: test_channel
	  username: test_user
	  webhook_url: test_webhook
	  threads: test_threads
  `
	reader, err := SubstituteConfigFromEnvVars(fname)
	require.Nil(t, err, "couldn't substitute config values: %s", err)

	bytes, err := io.ReadAll(reader)
	require.Nil(t, err, "couldn't read file data: %s", err)

	expectedFileContentLines := strings.Split(expectedFileContent, "\n")
	gotFileContentLines := strings.Split(string(bytes), "\n")
	for i := range expectedFileContentLines {
		require.Equal(t, expectedFileContentLines[i], gotFileContentLines[i], "lines in config don't match")
	}
}

func TestFileSizeToByteLen(t *testing.T) {
	byteLen, err := FileSizeToByteLen("2kb")
	require.Nil(t, err, "couldn't convert file size to byte len: %s", err)
	require.Equal(t, int(2048), byteLen)

	byteLen, err = FileSizeToByteLen("2mb")
	require.Nil(t, err, "couldn't convert file size to byte len: %s", err)
	require.Equal(t, int(2097152), byteLen)

	byteLen, err = FileSizeToByteLen("2")
	require.Nil(t, err, "couldn't convert file size to byte len: %s", err)
	require.Equal(t, int(2097152), byteLen)

	_, err = FileSizeToByteLen("2kilobytes")
	require.NotNil(t, err, "shouldn't convert file size to byte len: %s", err)
	require.ErrorContains(t, err, "parse error")
}

func TestOpenOrCreateFile(t *testing.T) {
	t.Run("should open an existing file", func(t *testing.T) {
		testFileName := "existingfile.txt"

		file, err := os.Create(testFileName)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		file.Close()

		file, err = OpenOrCreateFile(testFileName)
		require.NoError(t, err)
		require.True(t, FileExists(testFileName))
		file.Close()

		err = os.RemoveAll(testFileName)
		if err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	})

	t.Run("should create file if it does not exist", func(t *testing.T) {
		testFileName := "testfile.txt"
		file, err := OpenOrCreateFile(testFileName)
		require.NoError(t, err)
		require.True(t, FileExists(testFileName))
		file.Close()

		err = os.RemoveAll(testFileName)
		if err != nil {
			t.Fatalf("failed to remove test file: %v", err)
		}
	})

	t.Run("should fail when opening a non-existing file", func(t *testing.T) {
		testFileName := "/nonexistentdirectory/testfile.txt"
		_, err := OpenOrCreateFile(testFileName)

		require.Error(t, err)
	})
}

func TestFileExistsIn(t *testing.T) {
	tempDir := t.TempDir()
	anotherTempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "file.txt")
	err := os.WriteFile(tempFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	defer os.RemoveAll(tempFile)

	tests := []struct {
		name         string
		file         string
		allowedFiles []string
		expectedPath string
		expectedErr  bool
	}{
		{
			name:         "file exists in allowed directory",
			file:         tempFile,
			allowedFiles: []string{filepath.Join(tempDir, "tempfile.txt")},
			expectedPath: tempDir,
			expectedErr:  false,
		},
		{
			name:         "file does not exist in allowed directory",
			file:         tempFile,
			allowedFiles: []string{anotherTempDir},
			expectedPath: "",
			expectedErr:  true,
		},
		{
			name:         "path starting with .",
			file:         tempFile,
			allowedFiles: []string{"."},
			expectedPath: "",
			expectedErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			allowedPath, err := FileExistsIn(tc.file, tc.allowedFiles...)
			gotErr := err != nil
			require.Equal(t, tc.expectedErr, gotErr, "expected err but got %v", gotErr)
			require.Equal(t, tc.expectedPath, allowedPath)

		})
	}
}
