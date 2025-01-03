package dedupe

import (
	"prismx_cli/utils/hmap/store/hybrid"
)

type LevelDBBackend struct {
	storage *hybrid.HybridMap
}

func NewLevelDBBackend() *LevelDBBackend {
	l := &LevelDBBackend{}
	db, err := hybrid.New(hybrid.DefaultDiskOptions)
	if err != nil {
	}
	l.storage = db
	return l
}

func (l *LevelDBBackend) Upsert(elem string) bool {
	_, exists := l.storage.Get(elem)
	if exists {
		return false
	}

	if err := l.storage.Set(elem, nil); err != nil {
		return false
	}
	return true
}

func (l *LevelDBBackend) IterCallback(callback func(elem string)) {
	l.storage.Scan(func(k, _ []byte) error {
		callback(string(k))
		return nil
	})
}

func (l *LevelDBBackend) Cleanup() {
	_ = l.storage.Close()
}
