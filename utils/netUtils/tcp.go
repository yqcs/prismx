package netUtils

import (
	"github.com/saintfish/chardet"
	"net"
	"net/http"
	"net/http/httputil"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/cryptoPlus"
	"prismx_cli/utils/proxyPool"
	"prismx_cli/utils/randomUtils"
	"strings"
	"time"
)

// Pool URL格式 <schema>://<user>:<password>@<host>:<port>/<path>:<params>?<query>#<frag>
var openProxy bool

func OpenProxy(ipList []string) {
	openProxy = true
	proxyPool.IpList = arr.SliceRemoveDuplicates(ipList)
}

func CloseProxy() {
	openProxy = false
}

// Result 封装的http返回包
type Result struct {
	Other      *http.Response
	RequestRaw string
	Body       []byte
	Header     string
}

// SendHttp 自定义Http包
func SendHttp(request *http.Request, timeout time.Duration, redirect bool) (result Result, err error) {
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", randomUtils.GetUserAgent())
	}

	//获取请求Raw
	if requestOut, err := httputil.DumpRequestOut(request, true); err == nil {
		result.RequestRaw = string(requestOut)
	}

	result.Other, err = proxyPool.SendProxyHttp(openProxy, request, timeout, redirect)
	if err != nil {
		return result, err
	}

	if result.Other != nil {
		//无损取body
		result.Body = CopyRespBody(result.Other)
		//获取Header Raw
		headerOut, err := httputil.DumpResponse(result.Other, false)
		if err == nil {
			result.Header = string(headerOut)
		}
		var encoding string
		detectorStr, err := chardet.NewTextDetector().DetectBest(result.Body)
		if err != nil {
			encoding = cryptoPlus.GetEncoding(result.Other.Header.Get("Content-Type"), result.Body)
		} else {
			encoding = detectorStr.Charset
		}
		if strings.ToLower(encoding) != "utf-8" {
			result.Body = []byte(cryptoPlus.TransCode(result.Body, encoding))
		}
	}

	return result, err
}

func SendDialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return proxyPool.SendProxyTcp(openProxy, network, address, timeout)
}
