//go:build (darwin || linux) && (386 || arm)

package syscallutil

import "errors"

func loadLibrary(name string) (uintptr, error) {
	return 0, errors.New("not implemented")
}
