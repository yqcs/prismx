package owaspTop10

import (
	"github.com/gocolly/colly"
	"net/url"
	"prismx_cli/core/owaspTop10/fileIncloud"
	"prismx_cli/core/owaspTop10/sqli"
	"prismx_cli/core/owaspTop10/xss"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/randomUtils"
	"strings"
	"sync"
	"time"
)

type OwaspTop10 struct {
	Id, Target string
	Timeout    time.Duration
	State      string
}

func (t *OwaspTop10) Start() {

	var sqlBackUrl []string
	var fileBackUrl []string
	var xssBackUrl []string

	parse, err := url.Parse(t.Target)
	if err != nil {
		return
	}

	//JSFind函数已经存在此行为，可整合到一起，避免多次请求
	c := colly.NewCollector()
	c.UserAgent = randomUtils.GetUserAgent()
	c.AllowedDomains = []string{parse.Host}
	c.MaxDepth = 2

	//捕获form标签
	c.OnHTML("form[action]", func(e *colly.HTMLElement) {
		if t.State == "stop" {
			for {
				time.Sleep(3 * time.Second)
				if t.State == "run" {
					break
				}
				if t.State == "end" {
					return
				}
			}
		}
		if t.State == "end" {
			return
		}

		action := e.Attr("action")

		inputs := e.DOM.Find("input")
		var query string
		for i := 0; i < inputs.Length(); i++ {
			name, nameExists := inputs.Eq(i).Attr("name")
			if nameExists {
				value, valueExists := inputs.Eq(i).Attr("value")
				if valueExists {
					query += name + "=" + value + "&"
				} else {
					query += name + "=" + randomUtils.RandomString(5) + "&"
				}
			}
		}

		// Trim the trailing "&" from the query
		query = strings.TrimSuffix(query, "&")

		if query != "" {
			if !strings.HasPrefix(action, "/") {
				action = "/" + action
			}

			// Construct the full URL
			action = e.Request.URL.Path + action + "?" + query
		}
		e.Request.Visit(action)
	})
	//捕获a标签
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if t.State == "stop" {
			for {
				time.Sleep(3 * time.Second)
				if t.State == "run" {
					break
				}
				if t.State == "end" {
					return
				}
			}
		}
		if t.State == "end" {
			return
		}
		e.Request.Visit(e.Attr("href"))
	})
	//捕获img标签
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		if t.State == "stop" {
			for {
				time.Sleep(3 * time.Second)
				if t.State == "run" {
					break
				}
				if t.State == "end" {
					return
				}
			}
		}
		if t.State == "end" {
			return
		}
		e.Request.Visit(e.Attr("src"))
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.Scheme == "http" || r.URL.Scheme == "https" {
			logger.Info(r.URL.String())
			if t.State == "stop" {
				for {
					time.Sleep(3 * time.Second)
					if t.State == "run" {
						break
					}
					if t.State == "end" {
						return
					}
				}
			}
			if t.State == "end" {
				return
			}
			var wg sync.WaitGroup
			wg.Add(3)
			if t.State == "stop" {
				for {
					time.Sleep(3 * time.Second)
					if t.State == "run" {
						break
					}
					if t.State == "end" {
						return
					}
				}
			}

			go func() {
				if t.State == "run" {
					for _, item := range fileBackUrl {
						if item == r.URL.Path {
							wg.Done()
							return
						}
					}
					fileBackUrl = append(fileBackUrl, r.URL.Path)
					queryParams := r.URL.Query()

					// 遍历查询参数
					for key, values := range queryParams {
						for i := 0; i < len(values); i++ {
							// 检查是否包含等于号
							if values[i] != "" {
								queryParams.Set(key, "")
							}
						}
					}
					r.URL.RawQuery = queryParams.Encode()
					fileIncloud.Start(t.Id, r.URL, t.Timeout)
				}
				wg.Done()
			}()

			go func() {
				if t.State == "run" {
					for _, item := range sqlBackUrl {
						if item == r.URL.Path {
							wg.Done()
							return
						}
					}
					sqlBackUrl = append(sqlBackUrl, r.URL.Path)
					sqli.Start(t.Id, r.URL, t.Timeout)
				}
				wg.Done()
			}()

			go func() {
				if t.State == "run" {

					for _, item := range xssBackUrl {
						if item == r.URL.Path {
							wg.Done()
							return
						}
					}
					xssBackUrl = append(xssBackUrl, r.URL.Path)
					xss.Start(t.Id, r.URL, t.Timeout)
				}
				wg.Done()
			}()

			wg.Wait()
		}
	})
	c.Visit(t.Target)
}
