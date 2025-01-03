package proxyPool

import (
	"crypto/tls"
	"errors"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/logger"
	"strings"
	"time"
)

// proxyList 存活的代理 URL格式 <schema>://<user>:<password>@<host>:<port>/<path>:<params>?<query>#<frag>
var proxyList []string

// IpList 全部代理
var IpList []string

// RoundRobin 轮询 验证存活
// 检测存活代理  三分钟检查一次代理存活
func RoundRobin(timeout time.Duration) {
	for i := 0; i < len(IpList); i++ {
		u, err := url.Parse(IpList[i])
		if err != nil {
			continue
		}
		//校验是否存活
		conn, err := net.DialTimeout("tcp", u.Host, timeout)
		if err != nil {
			//不存活 删除
			for k := 0; k < len(proxyList); k++ {
				if proxyList[k] == IpList[i] {
					proxyList = append(proxyList[:k], proxyList[k+1:]...)
				}
			}
			continue
		}
		_ = conn.Close()
		//存活 加入池
		if !arr.IsContainByUrl(proxyList, IpList[i]) {
			proxyList = append(proxyList, IpList[i])
		}
	}
}

// SendProxyTcp 发送tcp代理
// 代理IP为空的话就不经过代理
func SendProxyTcp(proxyIP bool, network, address string, timeout time.Duration) (net.Conn, error) {
	var err error
	if proxyIP {
		if len(proxyList) == 0 {
			logger.Error("http: no available proxies")
			return nil, errors.New("no available proxies")
		}
		for i := 0; i < len(proxyList); i++ {
			u, er := url.Parse(proxyList[i])
			if er != nil {
				return nil, er
			}
			if strings.Contains(u.Scheme, "http") {
				logger.Error("Http proxy is not applicable to Tcp requests, and will directly use the direct connection mode")
				return net.DialTimeout(network, address, timeout)
			}
			auth := proxy.Auth{User: u.User.Username()}
			auth.Password, _ = u.User.Password()
			daile, er := proxy.SOCKS5("tcp", u.Host, &auth, &net.Dialer{Timeout: timeout})
			if er != nil {
				err = er
				continue
			}
			return daile.Dial(network, address)
		}
		return nil, err
	}
	return net.DialTimeout(network, address, timeout)
}

func SendProxyHttp(proxyIP bool, request *http.Request, timeout time.Duration, redirect bool) (*http.Response, error) {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10},
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: 10, //每个host最大空闲连接
		},
	}

	if !redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var err error
	if proxyIP {
		if len(proxyList) == 0 {
			logger.Error("tcp: no available proxies")
			return nil, errors.New("no available proxies")
		}
		for i := 0; i < len(proxyList); i++ {
			u, er := url.Parse(proxyList[i])
			if er != nil {
				continue
			}
			client.Transport = &http.Transport{
				Proxy:                 http.ProxyURL(u),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives:     true,
				MaxIdleConnsPerHost:   10,      //每个host最大空闲连接
				ResponseHeaderTimeout: timeout, //数据收发5秒超时
			}
			do, er := client.Do(request)
			if er != nil {
				err = er
				continue
			}
			return do, nil
		}
		return nil, err
	}
	return client.Do(request)
}
