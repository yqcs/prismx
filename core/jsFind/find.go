package jsFind

import (
	"crypto/tls"
	"embed"
	"github.com/gocolly/colly"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/url"
	"prismx_cli/utils/arr"
	"prismx_cli/utils/logger"
	"prismx_cli/utils/netUtils"
	"prismx_cli/utils/randomUtils"
	"regexp"
	"sync"
	"time"
)

type JsFind struct {
	Target, RuleName string
	Result           []string
}
type Rule struct {
	ID      string `yaml:"id"`
	Enabled bool   `yaml:"enabled"`
	Pattern string `yaml:"pattern"`
}

type Rules struct {
	Rule []Rule `yaml:"rules"`
}
type Machine struct {
	Target     string
	JsFindList chan JsFind
}

//go:embed rules.yaml
var rules embed.FS

func (m Machine) Start(timeout time.Duration) {

	data, err := rules.ReadFile("rules.yaml")
	if err != nil {
		logger.Error("dict file does not exist")
		return
	}

	var rulesList Rules

	if err = yaml.Unmarshal(data, &rulesList); err != nil {
		return
	}

	c := colly.NewCollector()
	c.UserAgent = randomUtils.GetUserAgent()

	parse, err := url.Parse(m.Target)
	if err != nil {
		return
	}

	c.AllowedDomains = []string{parse.Host}
	c.MaxDepth = 2

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("src"))
	})

	var jsList []string
	c.OnRequest(func(r *colly.Request) {
		jsList = append(jsList, r.URL.String())
	})

	if err = c.Visit(m.Target); err != nil {
		return
	}
	c.Wait()

	var wg sync.WaitGroup
	for _, js := range jsList {

		request, err := http.NewRequest("GET", js, nil)
		if err != nil {
			return
		}
		sendHttp, err := netUtils.SendHttp(request, timeout, false)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		for _, item := range rulesList.Rule {
			if !item.Enabled {
				continue
			}
			wg.Add(1)
			go func(rule Rule, u, body string) {
				defer wg.Done()
				reg := regexp.MustCompile(rule.Pattern)
				match := reg.FindAllString(body, -1)
				if len(match) > 0 {
					m.JsFindList <- JsFind{
						Target:   u,
						RuleName: rule.ID,
						Result:   arr.SliceRemoveDuplicates(match),
					}
				}
			}(item, js, string(sendHttp.Body))
		}
	}
	wg.Wait()
	close(m.JsFindList)
}
