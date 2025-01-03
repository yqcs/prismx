package scan

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	fingerprint "github.com/yqcs/fingerscan"
	"net"
	"prismx_cli/core/aliveCheck"
	"prismx_cli/core/hydra"
	"prismx_cli/core/models"
	"prismx_cli/core/plugins"
	"prismx_cli/core/subdomain/runner"
	"prismx_cli/core/vulnerability"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/parse"
	"prismx_cli/utils/task"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TaskInChan 下达任务前置操作
func (t *TaskPool) TaskInChan() {

	//处理域名，全部转小写
	for i := 0; i < len(t.Params.Uri); i++ {
		t.Params.Uri[i] = strings.ToLower(t.Params.Uri[i])
	}
	//将域名和IP放到一个列表
	hostList := append(t.Params.HostList, t.Params.Uri...)
	//去重
	hostList = arr.SliceRemoveDuplicates(hostList)
	//创建扫描任务
	hostScan := task.NewPool()
	//添加任务堵塞
	hostScan.Wg.Add(len(hostList))

	//并发执行函数，用来检测主机存活并下发端口扫描任务
	hostScan.PoolWithFunc, _ = ants.NewPoolWithFunc(t.Params.Thread, func(i interface{}) {

		defer hostScan.Wg.Done()

		host := i.(string)

		//判断是不是域名，并下发任务
		if !parse.IsDomain(host) {
			if t.Params.PN || aliveCheck.HostAliveCheck(host, t.Params.Ping, t.Params.Timeout, t.Params.AliveFuzz) {
				t.InvokeTask(host)
			}
			return
		}

		//扫描子域名
		if t.Params.SubDomain {
			t.Scan.Wg.Add(1)
			go func(domain string) {
				defer t.Scan.Wg.Done()

				var domainList []string
				//枚举子域名
				enumeration, _ := runner.RunEnumeration(runner.Runner{Target: domain, Timeout: t.Params.Timeout})
				//处理重复域名
				for _, item := range enumeration {
					if item.Value != "" && !strings.Contains(item.Value, "*") && !arr.IsContain(hostList, item.Value) && !arr.IsContain(domainList, item.Value) {
						domainList = append(domainList, item.Value)
					}
				}
				t.Scan.Wg.Add(len(domainList))
				for _, item := range domainList {
					go func(ipItem string) {
						defer t.Scan.Wg.Done()
						if t.Params.PN || aliveCheck.HostAliveCheck(ipItem, t.Params.Ping, t.Params.Timeout, t.Params.AliveFuzz) {
							//下发端口扫描任务
							t.InvokeTask(ipItem)
						}
					}(item)
				}
			}(host)
		}

		//扫描域名
		if t.Params.PN || aliveCheck.HostAliveCheck(host, t.Params.Ping, t.Params.Timeout, t.Params.AliveFuzz) {
			//下发端口扫描任务
			t.InvokeTask(host)
		}
	})

	//下发任务
	for _, item := range hostList {
		hostScan.PoolWithFunc.Invoke(item)
	}

	//实体队列堵塞
	hostScan.Wg.Wait()
	//清除任务
	hostScan.PoolWithFunc.Release()
	//实体队列堵塞
	t.Scan.Wg.Wait()
	//清除任务
	t.Scan.PoolWithFunc.Release()
}

// TaskFunc 端口扫描执行函数
func (t *TaskPool) TaskFunc(i any) {

	//--------------------扫描主机存活--------------------------

	//检测存活端口
	if !aliveCheck.PortCheck(i.(string), t.Params.Timeout) {
		return
	}
	host, port, err := net.SplitHostPort(i.(string))
	if err != nil {
		return
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return
	}

	//--------------------扫描指纹-------------------------

	//扫描指纹
	tcpFinger := fingerprint.ScanFingerprint(host, portInt, t.Params.Timeout+(5*time.Second))
	if tcpFinger == nil {
		return
	}
	//不留存响应包
	tcpFinger.Response = nil
	if tcpFinger.Version.VendorProductName != "unknown" {
		//移除重复
		tcpFinger.WebApp.App = arr.DeleteSliceValueToLower(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName)
		//如果检测到了version，将其拼接进appName里，组成 nginx 1.18.2
		if tcpFinger.Version.Version != "unknown" {
			tcpFinger.WebApp.App = append(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName+" "+tcpFinger.Version.Version)
		} else {
			tcpFinger.WebApp.App = append(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName)
		}
	}

	//格式化输出指纹数据
	appString := ""
	for _, item := range tcpFinger.WebApp.App {
		appString += logger.Global.Color().Green(" [") + logger.Global.Color().Cyan(item) + logger.Global.Color().Green("]")
	}
	if appString == "" {
		appString = " "
	}

	wsMsg := logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", tcpFinger.Uri)) + "	" + "[" + logger.Global.Color().Yellow(tcpFinger.Service) + "]" + "\t" + appString

	conMsg := net.JoinHostPort(host, port) + "\t" + tcpFinger.Service
	if tcpFinger.WebApp.App != nil {
		conMsg += " [" + strings.Join(tcpFinger.WebApp.App, ",") + "]"
	}

	if tcpFinger.WebApp.Title != "" {
		wsMsg += logger.Global.Color().Red("[") + tcpFinger.WebApp.Title + logger.Global.Color().Red("]")
		conMsg += " [" + tcpFinger.WebApp.Title + "]"
	}

	logger.ScanMessage(wsMsg)

	//检测漏洞，流程    ---先检测漏洞，然后再开始检测弱口令

	var wg sync.WaitGroup

	tcpFinger.Service = strings.TrimSuffix(tcpFinger.Service, "?")

	if strings.Contains(tcpFinger.Service, "ssl/http") {
		tcpFinger.Service = "https"
	}
	//检测漏洞
	if t.Params.Vul {
		wg.Add(2)

		//调用poc插件检测
		go func() {
			vulnerability.Verify(tcpFinger, t.Params.Timeout)
			wg.Done()
		}()

		go func() {
			//扫描top10
			//scan := owaspTop10.OwaspTop10{State: t.Scan.State, Id: t.ScanParams.TaskID, Target: tcpFinger.Service + "://" + tcpFinger.Uri, Timeout: t.ScanParams.Timeout}
			//scan.Start()
			wg.Done()
		}()

	}
	//根据协议匹配爆破组件
	if t.Params.WeakPass {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, item := range plugins.WeakPass {
				if tcpFinger.Service != item.App {
					continue
				}
				t.HydraTask = &models.HydraTask{
					Config: t.Params,
					App:    item.App,
					Target: net.JoinHostPort(tcpFinger.IP, strconv.Itoa(tcpFinger.Port)),
				}
				hydra.RunCheck(t.HydraTask, item)
			}
		}()
	}
	wg.Wait()
}

func (t *TaskPool) InvokeTask(s string) {
	for _, port := range t.Params.PortList {
		t.Scan.Wg.Add(1)
		t.Scan.PoolWithFunc.Invoke(net.JoinHostPort(s, strconv.Itoa(port)))
	}
}
