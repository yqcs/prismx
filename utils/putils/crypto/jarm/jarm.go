package jarm

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	gojarm "github.com/hdm/jarm-go"
	connpool "prismx_cli/utils/putils/conn/connpool"
)

// PoolCount defines how many connection are kept in the pool
var PoolCount = 3

// fingerprint probes a single host/port
func HashWithDialer(dialer connpool.Dialer, host string, port int, duration int) (string, error) {
	var results []string
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	timeout := time.Duration(duration) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(duration*PoolCount) * time.Second))
	defer cancel()

	// using connection pool as we need multiple probes
	pool, err := connpool.NewOneTimePool(ctx, addr, PoolCount)
	if err != nil {
		return "", err
	}
	pool.Dialer = dialer

	defer func() { _ = pool.Close() }()
	go func() { _ = pool.Run() }()

	for _, probe := range gojarm.GetProbes(host, port) {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			continue
		}
		if conn == nil {
			continue
		}
		_ = conn.SetWriteDeadline(time.Now().Add(timeout))
		_, err = conn.Write(gojarm.BuildProbe(probe))
		if err != nil {
			results = append(results, "")
			_ = conn.Close()
			continue
		}
		_ = conn.SetReadDeadline(time.Now().Add(timeout))
		buff := make([]byte, 1484)
		_, _ = conn.Read(buff)
		_ = conn.Close()
		ans, err := gojarm.ParseServerHello(buff, probe)
		if err != nil {
			results = append(results, "")
			continue
		}
		results = append(results, ans)
	}
	hash := gojarm.RawHashToFuzzyHash(strings.Join(results, ","))
	return hash, nil
}
