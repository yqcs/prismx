package utils

import (
	"net/url"
	"strings"
)

var WafKeyWord = []string{"造成安全威胁", "Bot-Block-ID", "您访问IP已被管理员限制", "本次事件ID", "当前访问疑似黑客攻击",
	"safedog", "拦截", "ValidateInputIfRequiredByConfig", "You don't have permission to access", "非法字符"}

func ParseQuery(target *url.URL, payload string) (item []string) {

	paramMap, err := url.ParseQuery(target.RawQuery)
	if err != nil {
		return
	}
	//如果没有抓到带参数的url，直接返回
	if len(paramMap) == 0 {
		return
	}
	for key, value := range paramMap {
		item = append(item, strings.Replace(target.String(), key+"="+value[0], key+"="+value[0]+payload, 1))
	}
	return item
}
