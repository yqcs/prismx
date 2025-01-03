package reverse

import (
	"prismx_cli/utils/interactsh/pkg/client"
	"prismx_cli/utils/interactsh/pkg/server"
	"prismx_cli/utils/randomUtils"
	"strings"
	"sync"
	"time"
)

type DnsServer struct {
	Client *client.Client
	Lock   *sync.RWMutex
}

var ResolveServer *DnsServer

var ResolveTaskMap = make(map[string]string)

// CheckServer 检查服务
func (d *DnsServer) CheckServer() string {
	if d == nil || d.Client == nil {
		return ""
	}
	return strings.ToLower(randomUtils.RandomString(10)) + "." + d.Client.URL()
}

func NewResolve() {
	c, err := client.New(client.DefaultOptions)
	if err != nil {
		ResolveServer = &DnsServer{Lock: &sync.RWMutex{}}
		return
	}
	ResolveServer = &DnsServer{Client: c, Lock: &sync.RWMutex{}}
	//运行反连监听程序
	go ResolveServer.Client.StartPolling(1*time.Second, func(interaction *server.Interaction) {
		ResolveServer.Lock.Lock()
		ResolveTaskMap[strings.ToLower(interaction.FullId)] = interaction.Protocol
		ResolveServer.Lock.Unlock()
		return
	})
}

// GetResolveUrl 检查数据
func GetResolveUrl() string {
	if ResolveServer == nil || ResolveServer.Client == nil {
		return ""
	}
	return strings.ToLower(randomUtils.RandomString(10)) + "." + ResolveServer.Client.URL()
}

// CheckResolveState 检查是否获取到数据
func CheckResolveState(protocol, url string, timeout time.Duration) bool {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return false
		default:
			for u, proto := range ResolveTaskMap {
				if proto == protocol && strings.HasPrefix(url, strings.ToLower(u)) {
					ResolveServer.Lock.Lock()
					delete(ResolveTaskMap, u)
					ResolveServer.Lock.Unlock()
					return true
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
