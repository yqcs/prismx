package plugins

import (
	"embed"
	"gopkg.in/yaml.v3"
	"os"
	"prismx_cli/core/hydra"
	"prismx_cli/core/models"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/file"
	"prismx_cli/utils/logger"
	"strings"
)

var (
	WeakPass []models.HydraAppFunc
	//go:embed exploits
	webApi embed.FS
)

func init() {

	//go reverse.NewResolve() //初始化dns

	WeakPass = []models.HydraAppFunc{
		{App: "vnc", Func: hydra.VncWeakPass},
		{App: "ftp", Func: hydra.FtpWeakPass},
		{App: "ms-wbt-server", Func: hydra.RdpWeakPass},
		{App: "wsman", Func: hydra.WinRMWeakPass},
		{App: "wsmans", Func: hydra.WinRMWeakPass},
		{App: "ssh", Func: hydra.SSHWeakPass},
		{App: "microsoft-ds", Func: hydra.SMBWeakPass},
		{App: "ms-sql-s", Func: hydra.MSSQLWeakPass},
		{App: "oracle-tns", Func: hydra.OracleWeakPass},
		{App: "mysql", Func: hydra.MySQLWeakPass},
		{App: "postgresql", Func: hydra.PGSQLWeakPass},
		{App: "redis", Func: hydra.RedisWeakPass},
		{App: "memcached", Func: hydra.MemcachedWeakPass},
		{App: "mongodb", Func: hydra.MongodbWeakPass},
		{App: "snmp", Func: hydra.SNMPWeakPass},
		{App: "zookeeper", Func: hydra.ZookeeperWeakPass},
		{App: "telnet", Func: hydra.TelnetWeakPass},
	}

	//加载yaml插件
	LoadAllYAML()

	//加载go插件
	LoadAllGoPoc()
}

// LoadAllYAML 加载全部插件
// 只能传插件目录
// 第一个参数是外部自定义的插件，第二个参数的内置poc路径
func LoadAllYAML() {

	models.YAMLPlugins = make(map[string][]models.AppVulInfo)

	//读取内置POC
	yamlList := file.FilesEmbedList(webApi, "exploits", "yaml")
	//读取json类型
	yamlList = append(yamlList, file.FilesEmbedList(webApi, "exploits", "json")...)

	//加载内置插件
	for _, item := range yamlList {
		result := LoadEmbedYAML(item)
		models.YAMLPlugins[result.App] = append(models.YAMLPlugins[result.App], result)
	}

	//读取外部自定义文件
	yamlList = file.FilesList("lib/exploits", "yaml")
	//读取json类型
	yamlList = append(yamlList, file.FilesList("lib/exploits", "json")...)
	//将外部自定义插件存入内存

	for _, item := range yamlList {
		result := LoadYAML(item)
		if result.App == "" {
			continue
		}
		//将App添加进去
		if len(models.YAMLPlugins[result.App]) == 0 {
			models.YAMLPlugins[result.App] = append(models.YAMLPlugins[result.App], result)
		}
		isFlag := true
		//检测当前插件是否存在
		for _, subItem := range models.YAMLPlugins[result.App] {
			if subItem.Meta.Name == result.Meta.Name {
				isFlag = false
				break
			}
		}
		if isFlag {
			models.YAMLPlugins[result.App] = append(models.YAMLPlugins[result.App], result)
		}
	}
}

// LoadAllGoPoc 加载go外部插件
func LoadAllGoPoc() {
	//读取外部自定义文件
	for _, item := range file.FilesList("lib/exploits", "go") {
		LoadGo(item)
	}
}

// LoadYAML 加载单个插件
func LoadYAML(s string) models.AppVulInfo {
	bytes, err := os.ReadFile(s)
	if err != nil {
		logger.Error(err.Error())
		return models.AppVulInfo{}
	}
	u := models.AppVulInfo{}
	if err = yaml.Unmarshal(bytes, &u); err != nil {
		logger.Error(err.Error())
	}
	return u
}

// LoadGo 加载单个Go插件
func LoadGo(s string) {

	//intp := interp.New(interp.Options{BuildTags: []string{"-s", "-w"}, GoPath: build.Default.GOPATH, Env: os.Environ()})
	//if err := intp.Use(unsafe.Symbols); err != nil {
	//	logger.Error(err.Error())
	//	return
	//}
	//if err := intp.Use(stdlib.Symbols); err != nil {
	//	logger.Error(err.Error())
	//	return
	//}
	//
	//err := intp.Use(map[string]map[string]reflect.Value{
	//	"prismx_cli/models/models": { //导入models
	//		"Register":      reflect.ValueOf(models.Register), //函数类型
	//		"HydraAppFunc":  reflect.ValueOf((*models.HydraAppFunc)(nil)),
	//		"VulResult":     reflect.ValueOf((*models.VulResult)(nil)),
	//		"AppVulInfo":    reflect.ValueOf((*models.AppVulInfo)(nil)), //models类型
	//		"VulMeta":       reflect.ValueOf((*models.VulMeta)(nil)),
	//		"StepsMeta":     reflect.ValueOf((*models.StepsMeta)(nil)),
	//		"VerifySteps":   reflect.ValueOf((*models.VerifySteps)(nil)),
	//		"ExploitSteps":  reflect.ValueOf((*models.ExploitSteps)(nil)),
	//		"ExploitParams": reflect.ValueOf((*models.ExploitParams)(nil)),
	//	},
	//	"prismx_cli/utils/netUtils/netUtils": { //导入netutils
	//		"SendHttp":        reflect.ValueOf(netUtils.SendHttp),
	//		"SendDialTimeout": reflect.ValueOf(netUtils.SendDialTimeout),
	//		"OpenProxy":       reflect.ValueOf(netUtils.OpenProxy),
	//		"CloseProxy":      reflect.ValueOf(netUtils.CloseProxy),
	//		"Result":          reflect.ValueOf((*netUtils.Result)(nil)),
	//	},
	//	"prismx_cli/utils/randomUtil/randomUtil": { //导出utils包
	//		"RandomString": reflect.ValueOf(randomUtils.RandomString),
	//		"GetUserAgent": reflect.ValueOf(randomUtils.GetUserAgent),
	//	},
	//	"prismx_cli/utils/reverse/reverse": { //导入反连包
	//		"CheckResolveState": reflect.ValueOf(reverse.CheckResolveState),
	//		"GetResolveUrl":     reflect.ValueOf(reverse.GetResolveUrl),
	//	},
	//})
	//if err != nil {
	//	logger.Error(err.Error())
	//	return
	//}
	//
	//if _, err = intp.EvalPath(s); err != nil {
	//	logger.Error(err.Error())
	//	return
	//}
}

// LoadEmbedYAML 加载单个内置插件
func LoadEmbedYAML(s string) models.AppVulInfo {
	readFile, err := webApi.ReadFile(s)
	if err != nil {
		logger.Error(err.Error())
		return models.AppVulInfo{}
	}
	u := models.AppVulInfo{}
	//u.Meta.Steps.VerifySteps.Verify = []models.StepMeta{}
	//u.Meta.Steps.ExploitSteps.Exploit = []models.StepsMate{}
	if err = yaml.Unmarshal(readFile, &u); err != nil {
		logger.Error(err.Error())
	}
	return u
}

// GetAllAppVulList 获取全部程序插件列表
func GetAllAppVulList() (list []Vul) {
	for _, item := range models.GOPlugins {
		for _, subItem := range item {
			list = append(list, Vul{"Service", subItem})
		}
	}
	for _, item := range models.YAMLPlugins {
		for _, subItem := range item {
			list = append(list, Vul{"WebApi", subItem})
		}
	}
	return list
}

// GetAllAppList 获取全部app名称
func GetAllAppList() (list []string) {
	for _, item := range models.GOPlugins {
		for _, subItem := range item {
			list = append(list, subItem.App)
		}
	}
	for _, item := range models.YAMLPlugins {
		for _, subItem := range item {
			list = append(list, subItem.App)
		}
	}
	return list
}

// GetGoAppList 获取全部Go程序
func GetGoAppList() (list []string) {
	for _, item := range models.GOPlugins {
		for _, subItem := range item {
			list = append(list, subItem.App)
		}
	}
	list = arr.SliceRemoveDuplicates(list)
	return list
}

// GetYAMLAppList 获取全部YAML程序
func GetYAMLAppList() (list []string) {
	for _, item := range models.YAMLPlugins {
		for _, subItem := range item {
			list = append(list, subItem.App)
		}
	}
	list = arr.SliceRemoveDuplicates(list)
	return list
}

// GetAllAppListByName 获取指定名称的app列表
func GetAllAppListByName(s string) (list []string) {
	for _, item := range GetGoAppListByAppName(s) {
		list = append(list, item)
	}
	for _, item := range GetYAMLAppListByAppName(s) {
		list = append(list, item)
	}
	return list
}

// GetAllAppListByNameSearch 模糊搜索App
func GetAllAppListByNameSearch(s string) (list []string) {
	//模糊匹配go app
	for k := range models.GOPlugins {
		if strings.Contains(strings.ToLower(k), strings.ToLower(s)) {
			list = append(list, k)
		}
	}

	//模糊匹配yml app
	for k := range models.YAMLPlugins {
		if strings.Contains(strings.ToLower(k), strings.ToLower(s)) {
			list = append(list, k)
		}
	}
	return list
}

// GetGoAppListByAppName 根据程序名称获取GOApp列表
func GetGoAppListByAppName(s string) (list []string) {
	for _, item := range models.GOPlugins[s] {
		list = append(list, item.App)
	}
	return list
}

// GetYAMLAppListByAppName 根据程序名称获取YAMLApp列表
func GetYAMLAppListByAppName(s string) (list []string) {
	for _, item := range models.YAMLPlugins[s] {
		list = append(list, item.App)
	}
	return list
}

// GetGoAppVulListByAppName 根据程序名称获取GO插件列表
func GetGoAppVulListByAppName(s string) (list []models.AppVulInfo) {
	for _, item := range models.GOPlugins[s] {
		list = append(list, item)
	}
	return list
}

// GetGoAppVulByVulName 根据漏洞名称获取go插件实体
func GetGoAppVulByVulName(s string) models.AppVulInfo {
	for _, item := range GetGoAppList() {
		for _, subItem := range models.GOPlugins[item] {
			if subItem.Meta.Name == s {
				return subItem
			}
		}
	}
	return models.AppVulInfo{}
}

// GetYAMLAppVulByVulName 根据漏洞名称获取yaml插件实体
func GetYAMLAppVulByVulName(s string) models.AppVulInfo {
	for _, item := range GetYAMLAppList() {
		for _, subItem := range models.YAMLPlugins[item] {
			if subItem.Meta.Name == s {
				return subItem
			}
		}
	}
	return models.AppVulInfo{}
}

// GetYAMLAppVulListByAppName 根据程序名称获取YAML插件列表
func GetYAMLAppVulListByAppName(s string) (list []models.AppVulInfo) {
	for _, item := range models.YAMLPlugins[s] {
		list = append(list, item)
	}
	return list
}

type Vul struct {
	PluginType string            `json:"plugin_type"`
	AppVul     models.AppVulInfo `json:"app_vul"`
}

// GetAppVulPluginListByAppName 根据App名称获取全部的插件
func GetAppVulPluginListByAppName(s string) (list []Vul) {
	for _, item := range GetGoAppVulListByAppName(s) {
		list = append(list, Vul{"Service", item})
	}
	for _, item := range GetYAMLAppVulListByAppName(s) {
		list = append(list, Vul{"WebApi", item})
	}
	return list
}

// GetAppVulListByAppName 根据App名称获取全部的插件，并且隐藏漏洞信息
func GetAppVulListByAppName(s string) (list []Vul) {
	for _, item := range GetGoAppVulListByAppName(s) {
		list = append(list, Vul{"Service", item})
	}
	for _, item := range GetYAMLAppVulListByAppName(s) {
		//给他隐藏poc
		item.Meta.Steps = models.VulMeta{}.Steps
		list = append(list, Vul{"WebApi", item})
	}
	return list
}
