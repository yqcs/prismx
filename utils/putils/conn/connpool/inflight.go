package connpool

import (
	"errors"
	"net"

	mapsutil "github.com/yqcs/prismx/utils/putils/maps"
	"go.uber.org/multierr"
)

type InFlightConns struct {
	inflightConns *mapsutil.SyncLockMap[net.Conn, struct{}]
}

func NewInFlightConns() (*InFlightConns, error) {
	m := &mapsutil.SyncLockMap[net.Conn, struct{}]{
		Map: mapsutil.Map[net.Conn, struct{}]{},
	}
	return &InFlightConns{inflightConns: m}, nil
}

func (i *InFlightConns) Add(conn net.Conn) {
	_ = i.inflightConns.Set(conn, struct{}{})
}

func (i *InFlightConns) Remove(conn net.Conn) {
	i.inflightConns.Delete(conn)
}

func (i *InFlightConns) Close() error {
	var errs []error

	_ = i.inflightConns.Iterate(func(conn net.Conn, _ struct{}) error {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
		return nil
	})

	if ok := i.inflightConns.Clear(); !ok {
		errs = append(errs, errors.New("couldn't empty in flight connections"))
	}

	return multierr.Combine(errs...)
}
