package main

import (
	"flag"
	"fmt"
	fingerprint "github.com/yqcs/fingerscan"
	"os"
	"prismx_cli/core/models"
	"prismx_cli/scan"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/parse"
	"prismx_cli/utils/reverse"
	"prismx_cli/utils/task"
	"strings"
	"time"
)

var (
	target    *string
	port      *string
	blackIP   *string
	blackPort *string
	mode      *string
	subdomain *bool
	ping      *bool
	weak      *bool
	vul       *bool
	pn        *bool
)

func init() {

	reverse.NewResolve()     //初始化dns
	fingerprint.InitFinger() //初始化指纹

	fmt.Println(`
  .
 /\\  ┌─┐  ╔═╗┬─┐┬┌─┐┌┬┐═╗ ╦
( ( ) │└┘  ╠═╝├┬┘│└─┐│││╔╩╦╝
 \\/  └──  ╩  ┴└─┴└─┘┴ ┴╩ ╚═
  '  :: PrismX Powerful security vulnerability scanning and attack threat capture tools
  //[Version OpenSource 1.0.0] (c) 2023 prismx.io All rights reserved. website[www.prismx.io]
`)

	//初始化账号密码
	password := strings.Split("123456,{user},,toor,qwe!2345,qwe!123,admin123,pass123,pass@123,password,123123,654321,111111,123,1,admin@123,Admin@123,admin123!@#,{user}1,{user}111,{user}123,{user}@123,{user}_123,{user}#123,{user}@111,{user}@2019,{user}@123#4,P@ssw0rd!,P@ssw0rd,Passw0rd,qwe123,12345678,test,test123,123qwe!@#,123456789,123321,666666,a123456.,123456~a,123456!a,000000,1234567890,8888888,!QAZ2wsx,1qaz2wsx,abc123,abc123456,1qaz@WSX,a11111,a12345,Aa1234,Aa1234.,Aa12345,a123456,a123123,Aa123123,Aa123456,Aa12345.,sysadmin,system,1qaz!QAZ,2wsx@WSX,qwe123!@#,Aa123456!,A123456s!,sa123456,1q2w3e", ",")

	models.UserDict["ftp"] = strings.Split("ftp,admin,www,web,root,db,wwwroot,data,anonymous", ",")
	models.PassDict["ftp"] = password
	models.UserDict["mysql"] = strings.Split("root,mysql", ",")
	models.PassDict["mysql"] = password
	models.UserDict["ms-sql-s"] = strings.Split("sql,sa", ",")
	models.PassDict["ms-sql-s"] = password
	models.UserDict["microsoft-ds"] = strings.Split("administrator,admin,guest", ",")
	models.PassDict["microsoft-ds"] = password
	models.UserDict["postgresql"] = strings.Split("postgres,admin,root", ",")
	models.PassDict["postgresql"] = password
	models.UserDict["ssh"] = strings.Split("root,admin", ",")
	models.PassDict["ssh"] = password
	models.UserDict["oracle"] = strings.Split("system,aqadm,sys,scott", ",")
	models.PassDict["oracle"] = password
	models.UserDict["redis"] = strings.Split("", ",")
	models.PassDict["redis"] = arr.DeleteSliceValue(password, "")
	models.UserDict["telnet"] = strings.Split("admin,root", ",")
	models.PassDict["telnet"] = password
	models.UserDict["wsman"] = strings.Split("administrator,admin", ",")
	models.PassDict["wsman"] = password
	models.UserDict["wsmans"] = strings.Split("administrator,admin", ",")
	models.PassDict["wsman"] = password
	models.UserDict["ms-wbt-server"] = strings.Split("administrator,admin,test", ",")
	models.PassDict["ms-wbt-server"] = password
	models.UserDict["memcached"] = strings.Split("", ",")
	models.PassDict["memcached"] = strings.Split("", ",")
	models.UserDict["mongodb"] = strings.Split("", ",")
	models.PassDict["mongodb"] = strings.Split("", ",")
	models.UserDict["zookeeper"] = strings.Split("", ",")
	models.PassDict["zookeeper"] = strings.Split("", ",")
	models.UserDict["snmp"] = strings.Split("", ",")
	models.PassDict["snmp"] = strings.Split("", ",")

	models.UserDict["vnc"] = strings.Split("", ",")
	models.PassDict["vnc"] = password

	target = flag.String("t", "", "Scan task target")
	subdomain = flag.Bool("s", false, "Scan subdomain, only applicable to the target containing domain")
	mode = flag.String("m", "d", "Scan speed: s, d, f")

	blackIP = flag.String("bip", "", "Filter host for scan task")
	blackPort = flag.String("bp", "", "Filter port for scan task")

	pn = flag.Bool("pn", false, "pn")
	vul = flag.Bool("vul", true, "Detect security risks")
	ping = flag.Bool("ping", false, "ICMP or Ping mode operation,default ICMP")
	weak = flag.Bool("weak", true, "Detect whether the common service has weak password")
	port = flag.String("p", "", "Ports to be scanned")
	flag.Parse()
}

func main() {

	if len(os.Args) == 1 || *target == "" {
		logger.Warn("Run command: -h")
		os.Exit(1)
	}
	//解析域名获得全部地址
	host, domainList, err := parse.ParseIP(*target, *blackIP)

	//域名列表和IP列表为空
	if err != nil {
		logger.Fatal(err)
	}

	params := models.ScanParams{
		Uri:       domainList,
		SubDomain: *subdomain,
		HostList:  host,
		IP:        *target,
		BlackIP:   *blackIP,
		BlackPort: *blackPort,
		Ping:      *ping,
		PN:        *pn,
		WeakPass:  *weak,
		Vul:       *vul,
		PortList:  parse.GetScanPort(*port, *blackPort),
	}

	if *mode == "s" {
		params.Thread = 800
		params.Timeout = 20 * time.Second
	}
	if *mode == "d" {
		params.Thread = 1500
		params.Timeout = 15 * time.Second
	}
	if *mode == "f" {
		params.Thread = 2500
		params.Timeout = 15 * time.Second
	}

	if len(params.PortList) == 0 {
		params.Port = "21,22,444,80,81,5040,4999,135,4630,139,443,445,1433,3306,5432,6379,8500,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,61616,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880"
		params.PortList = parse.GetScanPort(params.Port, params.BlackPort)
	}

	//创建任务实例
	pool := scan.TaskPool{Params: params, Scan: task.NewPool()}
	//启动任务
	pool.Start()

}
