package healthcheck

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type ConnectionInfo struct {
	Host       string
	Successful bool
	Message    string
	Error      error
}

func CheckConnection(host string, port int, protocol string, timeout time.Duration) ConnectionInfo {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout(protocol, address, timeout)
	if conn != nil {
		conn.Close()
	}

	return ConnectionInfo{
		Host:       host,
		Successful: err == nil,
		Message:    fmt.Sprintf("%s Connect (%s:%v): %s", protocol, host, port, "Successful"),
		Error:      err,
	}
}
