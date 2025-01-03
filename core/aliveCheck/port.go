package aliveCheck

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"
)

// PortCheck 尝试连接到指定地址的端口，返回是否成功
func PortCheck(addr string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func tcpScanPortCheck(ip string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ports := []int{80, 25, 22, 445, 8080, 8001, 8000, 9000, 23, 81, 6379, 21, 8500, 443, 3306, 8090}
	semaphore := make(chan struct{}, 10)
	var wg sync.WaitGroup
	openPorts := make(chan bool, len(ports))

	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			if PortCheck(net.JoinHostPort(ip, strconv.Itoa(port)), timeout) {
				openPorts <- true
			}
		}(port)
	}

	go func() {
		wg.Wait()
		close(openPorts)
	}()

	select {
	case <-ctx.Done():
		return false
	case res, ok := <-openPorts:
		if !ok {
			return res
		}
		return res
	}
}
