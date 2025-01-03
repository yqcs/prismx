package dedupe

// MaxInMemoryDedupeSize (default : 100 MB)
var MaxInMemoryDedupeSize = 100 * 1024 * 1024

type DedupeBackend interface {
	// Upsert add/update key to backend/database
	Upsert(elem string) bool
	// Execute given callback on each element while iterating
	IterCallback(callback func(elem string))
	// Cleanup cleans any residuals after deduping
	Cleanup()
}

// Dedupe is string deduplication type which removes
// all duplicates if
type Dedupe struct {
	receive <-chan string
	backend DedupeBackend
}

// Option is a type for variadic options in Drain
type Option func(val string)

// WithUnique is an option to send unique values to the provided channel
func WithUnique(ch chan<- string) Option {
	return func(val string) {
		ch <- val
	}
}

// Drains channel and tries to dedupe it
func (d *Dedupe) Drain(opts ...Option) {
	for val := range d.receive {
		if unique := d.backend.Upsert(val); unique {
			for _, opt := range opts {
				opt(val)
			}
		}
	}
}

// GetResults iterates over dedupe storage and returns results
func (d *Dedupe) GetResults() <-chan string {
	send := make(chan string, 100)
	go func() {
		defer close(send)
		d.backend.IterCallback(func(elem string) {
			send <- elem
		})
		d.backend.Cleanup()
	}()
	return send
}

// NewDedupe returns a dedupe instance which removes all duplicates
// Note: If byteLen is not correct/specified alterx may consume lot of memory
func NewDedupe(ch <-chan string, byteLen int) *Dedupe {
	d := &Dedupe{
		receive: ch,
	}
	if byteLen <= MaxInMemoryDedupeSize {
		d.backend = NewMapBackend()
	} else {
		// gologger print a info message here
		d.backend = NewLevelDBBackend()
	}
	return d
}
