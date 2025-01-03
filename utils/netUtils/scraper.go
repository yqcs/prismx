package netUtils

import (
	"net/http"
	"prismx_cli/utils/cryptoPlus"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type DocumentPreview struct {
	Title    string
	iconPath string
	Icon     string
}

func parseDocument(body []byte) (doc DocumentPreview) {

	// 初始化变量来保存标题和图标地址
	tokenizer := html.NewTokenizer(strings.NewReader(string(body)))

	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			return doc
		case html.SelfClosingTagToken, html.StartTagToken:
			token := tokenizer.Token()

			if token.Data == "link" {
				rel := ""
				href := ""
				for _, attr := range token.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					}
					if attr.Key == "href" {
						href = attr.Val
					}
				}
				if strings.Contains(rel, "icon") && doc.iconPath == "" {
					doc.iconPath = href
				}
			} else if token.Data == "title" {
				if tokenizer.Next() == html.TextToken {
					doc.Title = strings.ReplaceAll(tokenizer.Token().Data, "  ", "")
					doc.Title = strings.ReplaceAll(doc.Title, "\r", "")
					doc.Title = strings.ReplaceAll(doc.Title, "\n", "")
				}
			}
		}
	}
}

func Scrape(resp Result, timeout time.Duration) *DocumentPreview {
	doc := parseDocument(resp.Body)
	if doc == (DocumentPreview{}) {
		return nil
	}
	url := resp.Other.Request.URL.Scheme + "://" + resp.Other.Request.URL.Host + doc.iconPath
	//如果是独立地址
	if strings.HasPrefix(doc.iconPath, "http") {
		url = doc.iconPath
	} else if strings.HasPrefix(doc.iconPath, "//") {
		//如果是//exp.com地址
		url = resp.Other.Request.URL.Scheme + ":" + doc.iconPath
	} else if !strings.HasPrefix(doc.iconPath, "/") {
		// 如果是img.ico 则使用这种拼接
		url = resp.Other.Request.URL.Scheme + "://" + resp.Other.Request.URL.Host + "/" + doc.iconPath
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &doc
	}
	sendHttp, err := SendHttp(request, timeout, true)
	if err != nil {
		return &doc
	}
	if sendHttp.Other.Body != nil {
		sendHttp.Other.Body.Close()
	}
	if sendHttp.Other.StatusCode != http.StatusOK {
		return &doc
	}
	iconType := http.DetectContentType(sendHttp.Body)
	if strings.HasSuffix(doc.iconPath, ".svg") {
		doc.Icon = "data:image/svg+xml;base64," + cryptoPlus.Base64Encode(string(sendHttp.Body))
	} else if strings.Contains(iconType, "image/") {
		doc.Icon = "data:" + iconType + ";base64," + cryptoPlus.Base64Encode(string(sendHttp.Body))
	}
	return &doc
}

func GetTitle(body []byte) (title string) {
	if find := regexp.MustCompile("<title>(.*)</title>").FindSubmatch(body); find != nil {
		title = string(find[1])
	}
	return
}
