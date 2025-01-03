package syscallutil

func LoadLibrary(name string) (uintptr, error) {
	return loadLibrary(name)
}
