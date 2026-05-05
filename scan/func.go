package scan

import (
	"fmt"
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
	"strconv"
	"strings"
	"sync"
	"time"

	fingerprint "github.com/yqcs/fingerscan"
)

// TaskInChan 五阶段扫描流水线：
// 阶段1: 主机存活扫描 → 阶段2: 端口存活扫描 → 阶段3: 指纹识别 → 阶段4: 弱口令爆破 → 阶段5: 漏洞检测
func (t *TaskPool) TaskInChan() {
	// 处理域名，全部转小写
	for i := 0; i < len(t.Params.Uri); i++ {
		t.Params.Uri[i] = strings.ToLower(t.Params.Uri[i])
	}
	// 将域名和 IP 放到一个列表
	hostList := append(t.Params.HostList, t.Params.Uri...)
	// 去重
	hostList = arr.SliceRemoveDuplicates(hostList)

	var wg sync.WaitGroup

	// ==========================================
	// 阶段1: 主机存活扫描
	// ==========================================
	logger.Info(logger.Global.Color().Yellow(fmt.Sprintf("[阶段1/5] 开始主机存活扫描，共 %d 个目标", len(hostList))))

	var aliveHosts []string
	var aliveHostsMutex sync.Mutex

	hostChan := make(chan string, t.Params.Thread)

	for i := 0; i < t.Params.Thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range hostChan {
				// 判断是不是域名
				if !parse.IsDomain(host) {
					if t.Params.PN || aliveCheck.HostAliveCheck(host, t.Params.Ping, t.Params.Timeout, t.Params.AliveFuzz) {
						aliveHostsMutex.Lock()
						aliveHosts = append(aliveHosts, host)
						aliveHostsMutex.Unlock()
					}
					continue
				}

				// 扫描子域名
				if t.Params.SubDomain {
					var domainList []string
					enumeration, _ := runner.RunEnumeration(runner.Runner{Target: host, Timeout: t.Params.Timeout})
					for _, item := range enumeration {
						if item.Value != "" && !strings.Contains(item.Value, "*") && !arr.IsContain(hostList, item.Value) && !arr.IsContain(domainList, item.Value) {
							domainList = append(domainList, item.Value)
						}
					}
					aliveHostsMutex.Lock()
					aliveHosts = append(aliveHosts, domainList...)
					aliveHostsMutex.Unlock()
				}

				// 扫描域名本身
				if t.Params.PN || aliveCheck.HostAliveCheck(host, t.Params.Ping, t.Params.Timeout, t.Params.AliveFuzz) {
					aliveHostsMutex.Lock()
					aliveHosts = append(aliveHosts, host)
					aliveHostsMutex.Unlock()
				}
			}
		}()
	}

	// 发送任务
	for _, host := range hostList {
		hostChan <- host
	}
	close(hostChan)
	wg.Wait()

	// 去重存活主机列表
	aliveHosts = arr.SliceRemoveDuplicates(aliveHosts)
	logger.Info(logger.Global.Color().Green(fmt.Sprintf("[阶段1/5] 主机存活扫描完成，存活 %d 台", len(aliveHosts))))

	if len(aliveHosts) != 0 {
		for _, host := range aliveHosts {
			wsMsg := logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", host)) + "	" + "[" + logger.Global.Color().Yellow("Alive") + "]"
			logger.ScanMessage(wsMsg)
		}
	} else {
		logger.Info(logger.Global.Color().Yellow("没有存活主机，扫描结束"))
		return
	}

	// ==========================================
	// 阶段2: 端口存活扫描
	// ==========================================
	var totalPorts = len(aliveHosts) * len(t.Params.PortList)
	logger.Info(logger.Global.Color().Yellow(fmt.Sprintf("[阶段2/5] 开始端口存活扫描，共 %d 个端口任务", totalPorts)))

	var (
		alivePorts      []string
		alivePortsMutex sync.Mutex
	)

	portChan := make(chan string, t.Params.Thread)

	for i := 0; i < t.Params.Thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range portChan {
				if aliveCheck.PortCheck(target, t.Params.Timeout) {
					wsMsg := logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", target)) + "	" + "[" + logger.Global.Color().Yellow("Open") + "]"
					logger.ScanMessage(wsMsg)
					alivePortsMutex.Lock()
					alivePorts = append(alivePorts, target)
					alivePortsMutex.Unlock()
				}
			}
		}()
	}

	// 发送端口扫描任务
	for _, host := range aliveHosts {
		for _, port := range t.Params.PortList {
			portChan <- net.JoinHostPort(host, strconv.Itoa(port))
		}
	}
	close(portChan)
	wg.Wait()

	logger.Info(logger.Global.Color().Green(fmt.Sprintf("[阶段2/5] 端口存活扫描完成，存活端口 %d 个", len(alivePorts))))

	if len(alivePorts) == 0 {
		logger.Info(logger.Global.Color().Yellow("没有存活端口，扫描结束"))
		return
	}

	// ==========================================
	// 阶段3: 指纹识别
	// ==========================================
	logger.Info(logger.Global.Color().Yellow(fmt.Sprintf("[阶段3/5] 开始指纹识别，共 %d 个目标", len(alivePorts))))

	var (
		fingerResults      []*fingerprint.AppFinger
		fingerResultsMutex sync.Mutex
	)

	fingerChan := make(chan string, t.Params.Thread)

	for i := 0; i < t.Params.Thread; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range fingerChan {
				host, port, err := net.SplitHostPort(target)
				if err != nil {
					continue
				}
				portInt, err := strconv.Atoi(port)
				if err != nil {
					continue
				}

				// 扫描指纹
				tcpFinger := fingerprint.ScanFingerprint(host, portInt, t.Params.Timeout+(5*time.Second))
				if tcpFinger == nil {
					continue
				}

				// 不留存响应包
				tcpFinger.Response = nil
				if tcpFinger.Version.VendorProductName != "unknown" {
					// 移除重复
					tcpFinger.WebApp.App = arr.DeleteSliceValueToLower(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName)
					// 如果检测到了version，将其拼接进appName里，组成 nginx 1.18.2
					if tcpFinger.Version.Version != "unknown" {
						tcpFinger.WebApp.App = append(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName+" "+tcpFinger.Version.Version)
					} else {
						tcpFinger.WebApp.App = append(tcpFinger.WebApp.App, tcpFinger.Version.VendorProductName)
					}
				}

				// 格式化输出指纹数据
				appString := ""
				for _, item := range tcpFinger.WebApp.App {
					appString += logger.Global.Color().Green(" [") + logger.Global.Color().Cyan(item) + logger.Global.Color().Green("]")
				}
				if appString == "" {
					appString = " "
				}

				wsMsg := logger.Global.Color().Green(fmt.Sprintf("%-"+strconv.Itoa(20)+"s", tcpFinger.Uri)) + "\t" + "[" + logger.Global.Color().Yellow(tcpFinger.Service) + "]" + "\t" + appString

				if tcpFinger.WebApp.Title != "" {
					wsMsg += logger.Global.Color().Red("[") + tcpFinger.WebApp.Title + logger.Global.Color().Red("]")
				}

				logger.ScanMessage(wsMsg)

				fingerResultsMutex.Lock()
				fingerResults = append(fingerResults, tcpFinger)
				fingerResultsMutex.Unlock()
			}
		}()
	}

	// 发送指纹识别任务
	for _, target := range alivePorts {
		fingerChan <- target
	}
	close(fingerChan)
	wg.Wait()

	logger.Info(logger.Global.Color().Green(fmt.Sprintf("[阶段3/5] 指纹识别完成，成功识别 %d 个服务", len(fingerResults))))

	if len(fingerResults) == 0 {
		logger.Info(logger.Global.Color().Yellow("没有识别到服务指纹，扫描结束"))
		return
	}

	// ==========================================
	// 阶段4: 弱口令爆破
	// ==========================================
	if t.Params.WeakPass {
		logger.Info(logger.Global.Color().Yellow(fmt.Sprintf("[阶段4/5] 开始弱口令爆破，共 %d 个目标", len(fingerResults))))

		weakChan := make(chan *fingerprint.AppFinger, t.Params.Thread)

		for i := 0; i < t.Params.Thread; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for tcpFinger := range weakChan {
					// 规范化服务名
					service := strings.TrimSuffix(tcpFinger.Service, "?")
					if strings.Contains(service, "ssl/http") {
						service = "https"
					}

					for _, item := range plugins.WeakPass {
						if service != item.App {
							continue
						}
						hydraTask := &models.HydraTask{
							Config: t.Params,
							App:    item.App,
							Target: net.JoinHostPort(tcpFinger.IP, strconv.Itoa(tcpFinger.Port)),
						}
						hydra.RunCheck(hydraTask, item)
					}
				}
			}()
		}

		// 发送弱口令爆破任务
		for _, result := range fingerResults {
			weakChan <- result
		}
		close(weakChan)
		wg.Wait()

		logger.Info(logger.Global.Color().Green("[阶段4/5] 弱口令爆破完成"))
	} else {
		logger.Info(logger.Global.Color().Yellow("[阶段4/5] 弱口令爆破已跳过"))
	}

	// ==========================================
	// 阶段5: 漏洞检测
	// ==========================================
	if t.Params.Vul {
		logger.Info(logger.Global.Color().Yellow(fmt.Sprintf("[阶段5/5] 开始漏洞检测，共 %d 个目标", len(fingerResults))))

		vulnChan := make(chan *fingerprint.AppFinger, t.Params.Thread)

		for i := 0; i < t.Params.Thread; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for tcpFinger := range vulnChan {
					vulnerability.Verify(tcpFinger, t.Params.Timeout)
				}
			}()
		}

		// 发送漏洞检测任务
		for _, result := range fingerResults {
			vulnChan <- result
		}
		close(vulnChan)
		wg.Wait()

		logger.Info(logger.Global.Color().Green("[阶段5/5] 漏洞检测完成"))
	} else {
		logger.Info(logger.Global.Color().Yellow("[阶段5/5] 漏洞检测已跳过"))
	}

	logger.Info(logger.Global.Color().Green("全部扫描阶段完成"))
}
