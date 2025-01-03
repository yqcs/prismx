//go:build windows

package syscallutil

import "golang.org/x/sys/windows"

func loadLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadLibrary(name)
	return uintptr(handle), err
}
