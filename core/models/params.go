package models

import (
	"prismx_cli/utils/task"
	"time"
)

var PassDict = map[string][]string{}

var UserDict = map[string][]string{}

// Dict 组合密码字典
type Dict struct {
	User     string
	Password string
}

// HydraTask 爆破任务配置
type HydraTask struct {
	App      string
	Target   string
	DictList []Dict
	Config   ScanParams
	Dict     Dict
	Scan     *task.Pool
}

type ScanParams struct {
	IP, Port, BlackPort, BlackIP                                      string        //要检查的IP、端口以及黑名单
	Timeout                                                           time.Duration //网络超时
	Thread                                                            int           //线程数
	Ping, AliveCheck, AliveFuzz, WeakPass, SubDomain, PN, Vul, Nuclei bool          //是否用ping、智能/模糊模式检测存活、基线扫描、目录扫描、子域名扫描、PN、vul扫描
	HostList                                                          []string      //解析得到的待扫描IP
	PortList                                                          []int         //解析得到的待扫描端口
	Uri                                                               []string      //域名列表
}
