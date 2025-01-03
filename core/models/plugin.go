package models

import (
	"time"
)

var (
	YAMLPlugins = make(map[string][]AppVulInfo) //yaml插件集
	GOPlugins   = make(map[string][]AppVulInfo) //go所有插件
)

// HydraAppFunc 爆破使用
type HydraAppFunc struct {
	App  string        //协议
	Func func(any) any //爆破函数
}

// Register 注册插件
func Register(Info AppVulInfo) {
	for _, item := range GOPlugins[Info.App] {
		if item.Meta.Name == Info.Meta.Name {
			return
		}
	}
	GOPlugins[Info.App] = append(GOPlugins[Info.App], Info)
}

// VulResult 封装漏洞检测返回信息
type VulResult struct {
	Request  string `json:"request,omitempty"`
	Response string `json:"response"`
	State    bool   `json:"state"`
}

// AppVulInfo 程序漏洞信息
type AppVulInfo struct {
	App   string  `json:"app" form:"app" yaml:"app"` //App名称
	Query string  `json:"query" yaml:"query"`        //程序类型 1 基线 2 web
	Meta  VulMeta `json:"meta" yaml:"meta" gorm:"-"` //漏洞信息
}

// VulMeta 每个漏洞详细信息
type VulMeta struct {
	Name        string    `json:"name" yaml:"name"`               //漏洞名称
	Level       int       `json:"level" yaml:"level"`             //危险级别 信息、低、中、高、严重 Nil  = 0 Info = 1 Low = 2 Middle = 3 High = 4、Serious = 5
	Tags        []string  `json:"tags" yaml:"tags"`               //漏洞类型，如远程命令执行
	Description string    `json:"description" yaml:"description"` //漏洞描述
	Homepage    string    `json:"homepage" yaml:"homepage"`       //厂商官网
	Author      string    `json:"author" yaml:"author"`           //漏洞作者
	References  string    `json:"references" yaml:"references"`   //参考文章
	Solution    string    `json:"solution" yaml:"solution"`       //解决方案
	CreateAt    string    `json:"create_at" yaml:"create_at"`     //poc创建时间
	Available   bool      `json:"available" yaml:"available"`     //可利用
	Steps       StepsMeta `json:"steps" yaml:"steps"`
}

type StepsMeta struct {
	Variable     []map[string]any `json:"variable" yaml:"variable"`                               //每个请求都允许自定义 允许存在request：header、path、params、body response：code、body、raw中
	VerifySteps  VerifySteps      `json:"verify_steps" yaml:"verify_steps"`                       //POC模块，漏洞验证
	ExploitSteps ExploitSteps     `json:"exploit_steps,omitempty" yaml:"exploit_steps,omitempty"` //exp模块
}

type VerifySteps struct {
	Type     string                                                              `json:"type" yaml:"type"`                         //or | and  可能有多种复现方式。这个时候就使用此参数决定响应内容是都必须符合还是只需符合其中一项
	Verify   []StepMeta                                                          `json:"verify,omitempty" yaml:"verify,omitempty"` //poc验证，如果是数组的情况，type为and情况下最后一次匹配到的结果为最终请求与响应数据，or的情况下第一次匹配规则即是最终响应。
	VerifyGo func(scheme, ip string, port int, duration time.Duration) VulResult `json:"-" yaml:"-"`
}

type ExploitSteps struct {
	Type      string                                                                              `json:"type" yaml:"type"` //or | and  可能有多种复现方式。这个时候就使用此参数决定响应内容是都必须符合还是只需符合其中一项
	Params    ExploitParams                                                                       `json:"params" yaml:"params"`
	Exploit   []StepMeta                                                                          `json:"exploit,omitempty" yaml:"exploit,omitempty"` //exp验证，如果是数组的情况，type为and情况下最后一次匹配到的结果为最终请求与响应数据，or的情况下第一次匹配规则即是最终响应。
	ExploitGo func(scheme, ip string, port int, payload string, duration time.Duration) VulResult `json:"-" yaml:"-"`
}

type ExploitParams struct {
	Name  string `json:"name" yaml:"name"`   //选择框标题，如 CMD、FILE
	Type  string `json:"type" yaml:"type"`   //类型，支持 input、select
	Value string `json:"value" yaml:"value"` //select多参数需以逗号分割  select： key:value  如 readFile:/index?path=/etc/passwd  预留占位符：{{exploit}}
}

// StepMeta 验证结构体
type StepMeta struct {
	Request struct {
		Method   string              `json:"method" yaml:"method"`                    //请求类型
		Path     string              `json:"path" yaml:"path"`                        //存在问题的地址
		Redirect bool                `json:"redirect" yaml:"redirect" default:"true"` //http是否允许重定向，默认允许
		Header   []map[string]string `json:"header" yaml:"header"`
		Params   string              `json:"params" yaml:"params"` //请求参数
	} `json:"request" yaml:"request"` //构造payload
	Response []struct {
		Name  string `json:"name" yaml:"name"`   //校验响应的数据模块，如header、body、code、time
		Value string `json:"value" yaml:"value"` //期待的值
		Type  string `json:"type" yaml:"type"`   //操作类型，如等于=、包含contain、不等于!=、不包含not contain、regex
	} `json:"response" yaml:"response"` //根据响应信息检测漏洞是否存在 示例： code == 200
}

// MSG success out payload
type MSG struct {
	Target    string        `json:"target"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	Payload   any           `json:"payload"`
	Response  string        `json:"response"`
	Leve      int           `json:"leve"`
	Describe  string        `json:"describe"`
	EXP       bool          `json:"exp"`
	ExpParams ExploitParams `json:"exp_params"`
}
