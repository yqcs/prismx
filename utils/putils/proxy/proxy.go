package proxyutils

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	errorutil "prismx_cli/utils/putils/errors"
)

type proxyResult struct {
	AliveProxy string
	Error      error
}

// ProxyProbeConcurrency (default 8)
var ProxyProbeConcurrency = 8

const (
	SOCKS5 = "socks5"
	HTTP   = "http"
	HTTPS  = "https"
)

// GetAnyAliveProxy takes proxies as input and returns the first alive proxy
// or returns error if all of them not alive
func GetAnyAliveProxy(timeoutInSec int, proxies ...string) (string, error) {
	sg := sizedwaitgroup.New(ProxyProbeConcurrency)
	resChan := make(chan proxyResult, 4)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for _, v := range proxies {
			// skip iterating if alive proxy is found
			select {
			case <-ctx.Done():
				return
			default:
				proxy, err := GetProxyURL(v)
				if err != nil {
					resChan <- proxyResult{Error: err}
					continue
				}
				sg.Add()
				go func(proxyAddr url.URL) {
					defer sg.Done()
					select {
					case <-ctx.Done():
						return
					case resChan <- testProxyConn(proxyAddr, timeoutInSec):
						cancel()
					}
				}(proxy)
			}
		}
		sg.Wait()
		close(resChan)
	}()

	errstack := []string{}
	for {
		result, ok := <-resChan
		if !ok {
			break
		}
		if result.AliveProxy != "" {
			// found alive proxy return now
			return result.AliveProxy, nil
		} else if result.Error != nil {
			errstack = append(errstack, result.Error.Error())
		}
	}

	// all proxies are dead
	return "", errorutil.NewWithTag("proxyutils", "all proxies are dead got : %v", strings.Join(errstack, " : "))
}

// dial and test if proxy is open
func testProxyConn(proxyAddr url.URL, timeoutInSec int) proxyResult {
	p := proxyResult{}
	if Conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", proxyAddr.Hostname(), proxyAddr.Port()), time.Duration(timeoutInSec)*time.Second); err == nil {
		_ = Conn.Close()
		p.AliveProxy = proxyAddr.String()
	} else {
		p.Error = err
	}
	return p
}

// GetProxyURL returns a Proxy URL after validating if given proxy url is valid
func GetProxyURL(proxyAddr string) (url.URL, error) {
	if url, err := url.Parse(proxyAddr); err == nil && isSupportedProtocol(url.Scheme) {
		return *url, nil
	}
	return url.URL{}, errorutil.New("invalid proxy format (It should be http[s]/socks5://[username:password@]host:port)").WithTag("proxyutils")
}

// isSupportedProtocol checks given protocols are supported
func isSupportedProtocol(value string) bool {
	return value == HTTP || value == HTTPS || value == SOCKS5
}
