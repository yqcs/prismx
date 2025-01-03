package channelutil

// CreateNChannels creates and returns N channels
func CreateNChannels[T any](count int, bufflen int) map[int]chan T {
	x := map[int]chan T{}

	for i := 0; i < count; i++ {
		x[i] = make(chan T, bufflen)
	}
	return x
}
