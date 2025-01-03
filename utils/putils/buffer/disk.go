package buffer

import (
	"io"
	"os"
)

type DiskBuffer struct {
	f *os.File
}

func New() (*DiskBuffer, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}

	return &DiskBuffer{f: f}, nil
}

func (db *DiskBuffer) Write(b []byte) (int, error) {
	return db.f.Write(b)
}

func (db *DiskBuffer) WriteAt(b []byte, off int64) (int, error) {
	return db.f.WriteAt(b, off)
}

func (db *DiskBuffer) WriteString(s string) (int, error) {
	return db.f.WriteString(s)
}

func (db *DiskBuffer) Bytes() ([]byte, error) {
	return os.ReadFile(db.f.Name())
}

func (db *DiskBuffer) String() (string, error) {
	data, err := db.Bytes()
	return string(data), err
}

// all readers must be closed to avoid FD leak
func (db *DiskBuffer) Reader() (io.ReadSeekCloser, error) {
	f, err := os.Open(db.f.Name())
	return f, err
}

func (db *DiskBuffer) Close() {
	name := db.f.Name()
	db.f.Close()
	os.RemoveAll(name)
}
