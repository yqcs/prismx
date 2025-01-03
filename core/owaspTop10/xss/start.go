package xss

import (
	"net/http"
	"net/url"
	"prismx_cli/core"
	"prismx_cli/core/models"
	"prismx_cli/core/owaspTop10/utils"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/netUtils"
	"prismx_cli/utils/randomUtils"
	"strconv"
	"strings"
	"time"
)

func Start(id string, target *url.URL, timeout time.Duration) {
	rand := randomUtils.RandomString(5)
	payload := []string{`<b>` + rand + `</b>`, `<zzz><ScRiPt>` + rand + `</ScRiPt>`}
	for _, item := range payload {
		for _, subItem := range utils.ParseQuery(target, item) {
			go func(uri string) {
				request, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					return
				}
				sendHttp, err := netUtils.SendHttp(request, timeout, false)
				if err != nil || sendHttp.Other.StatusCode == 404 {
					return
				}
				//先检查waf
				if arr.IsContain(utils.WafKeyWord, string(sendHttp.Body)) {
					return
				}
				port, _ := strconv.Atoi(target.Port())
				if strings.Contains(string(sendHttp.Body), item) && !strings.Contains(string(sendHttp.Body), "Warning: fopen(") {
					core.CaptureVul(id, target.Hostname(), port, models.MSG{
						Name:     "Cross Site Scripting",
						Type:     "XSS",
						Payload:  sendHttp.RequestRaw,
						Target:   target.Scheme + "://" + target.Host,
						Response: sendHttp.Header + string(sendHttp.Body),
						Describe: "Cross site scripting attack XSS involves injecting malicious script code into a web page. When a user browses the page, the script code embedded in the web will be executed, thereby achieving the goal of malicious attack on the user.",
						Leve:     3,
						EXP:      false,
					})
					return
				}
			}(subItem)
		}
	}
}
