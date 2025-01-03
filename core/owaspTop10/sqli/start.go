package sqli

import (
	"embed"
	"github.com/beevik/etree"
	"io"
	"net/http"
	"net/url"
	"prismx_cli/core"
	"prismx_cli/core/models"
	"prismx_cli/core/owaspTop10/utils"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/netUtils"
	"prismx_cli/utils/randomUtils"
	"regexp"
	"strconv"
	"time"
)

//go:embed data
var dataFs embed.FS

func initSqlInject() map[string]string {
	var errorsXml = make(map[string]string)
	errors, err := dataFs.Open("data/errors.xml")
	if err != nil {
		return nil
	}
	bytes, _ := io.ReadAll(errors)
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(bytes); err != nil {
		panic(err)
	}
	for _, dbms := range doc.SelectElement("root").SelectElements("dbms") {
		for _, item := range dbms.Attr {
			for _, e := range dbms.SelectElements("error") {
				for _, subItem := range e.Attr {
					errorsXml[subItem.Value] = item.Value
				}
			}
		}
	}
	return errorsXml
}

func Start(id string, target *url.URL, timeout time.Duration) {
	//先测试单引号爆错，再测试and逻辑运算符
	payload := []string{"1'" + randomUtils.RandomString(5), "1 and '" + randomUtils.RandomString(4) + "' = '" + randomUtils.RandomString(5) + "'"}

	//报错类型poc
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
				if arr.IsContainToLower(utils.WafKeyWord, string(sendHttp.Body)) {
					return
				}

				//检查是否包含错误关键词
				for subKey, subValue := range initSqlInject() {
					match, _ := regexp.MatchString(subKey, string(sendHttp.Body))
					if match == true {
						port, _ := strconv.Atoi(target.Port())
						core.CaptureVul(id, target.Hostname(), port, models.MSG{
							Name:     subValue + " Inject",
							Type:     "SQLI",
							Payload:  sendHttp.RequestRaw,
							Target:   target.Scheme + "://" + target.Host,
							Response: sendHttp.Header + string(sendHttp.Body),
							Describe: "SQL injection vulnerability refers to an attacker inserting malicious SQL statements into website parameters through a browser or other client, but the website application does not filter them, bringing malicious SQL statements into the database for execution, thereby allowing the attacker to obtain sensitive information or perform other malicious operations through the database.",
							Leve:     4,
							EXP:      false,
						})
						return
					}
				}
			}(subItem)

		}
	}
}
