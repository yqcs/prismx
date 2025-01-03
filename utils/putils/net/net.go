package netutil

import (
	"errors"
	"net"
)

var ErrMissingPort = errors.New("missing port")

// TryJoinHostPort joins host and port. If port is empty, it returns host and an error.
func TryJoinHostPort(host, port string) (string, error) {
	if host == "" {
		return "", &net.AddrError{Err: "missing host", Addr: host}
	}

	if port == "" {
		return host, ErrMissingPort
	}

	return net.JoinHostPort(host, port), nil
}
