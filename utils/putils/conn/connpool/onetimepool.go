package connpool

import (
	"context"
	"net"
	"sync"
)

type Dialer interface {
	Dial(ctx context.Context, network, address string) (net.Conn, error)
}

// OneTimePool is a pool designed to create continous bare connections that are for one time only usage
type OneTimePool struct {
	address         string
	idleConnections chan net.Conn
	InFlightConns   *InFlightConns
	ctx             context.Context
	cancel          context.CancelFunc
	Dialer          Dialer
	mx              sync.RWMutex
}

func NewOneTimePool(ctx context.Context, address string, poolSize int) (*OneTimePool, error) {
	idleConnections := make(chan net.Conn, poolSize)
	inFlightConns, err := NewInFlightConns()
	if err != nil {
		return nil, err
	}
	pool := &OneTimePool{
		address:         address,
		idleConnections: idleConnections,
		InFlightConns:   inFlightConns,
	}
	if ctx == nil {
		ctx = context.Background()
	}
	pool.ctx, pool.cancel = context.WithCancel(ctx)
	return pool, nil
}

// Acquire acquires an idle connection from the pool
func (p *OneTimePool) Acquire(c context.Context) (net.Conn, error) {
	select {
	case <-p.ctx.Done():
		return nil, p.ctx.Err()
	case <-c.Done():
		return nil, c.Err()
	case conn := <-p.idleConnections:
		p.InFlightConns.Remove(conn)
		return conn, nil
	}
}

func (p *OneTimePool) Run() error {
	for {
		select {
		case <-p.ctx.Done():
			return p.ctx.Err()
		default:
			var (
				conn net.Conn
				err  error
			)
			p.mx.RLock()
			hasDialer := p.Dialer != nil
			p.mx.RUnlock()

			if hasDialer {
				p.mx.RLock()
				conn, err = p.Dialer.Dial(p.ctx, "tcp", p.address)
				p.mx.RUnlock()
			} else {
				conn, err = net.Dial("tcp", p.address)
			}
			if err == nil {
				p.InFlightConns.Add(conn)
				p.idleConnections <- conn
			}
		}
	}
}

func (p *OneTimePool) Close() error {
	p.cancel()

	// remove dialer references
	p.mx.Lock()
	p.Dialer = nil
	p.mx.Unlock()

	return p.InFlightConns.Close()
}
