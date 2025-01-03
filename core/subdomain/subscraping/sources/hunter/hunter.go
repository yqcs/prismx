package hunter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
	"prismx_cli/utils/cryptoPlus"
)

type hunterResponse struct {
	Code int `json:"code"`
	Data struct {
		Total int `json:"total"`
		Time  int `json:"time"`
		Arr   []struct {
			Url      string `json:"url"`
			IP       string `json:"ip"`
			Port     int    `json:"port"`
			Domain   string `json:"domain"`
			WebTitle string `json:"web_title"`
		} `json:"arr"`
	} `json:"data"`
	Message string `json:"message"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.HunterUserName == "" || session.Keys.HunterKey == "" {
			return
		}

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://hunter.qianxin.com/openApi/search?username=%s&api-key=%s&search=%s&page=1&page_size=10&is_web=1", session.Keys.HunterUserName, session.Keys.HunterKey, cryptoPlus.Base64Encode(domain)))

		if err != nil && resp == nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		defer resp.Body.Close()
		var response hunterResponse
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &response); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if response.Code != 200 {
			results <- subscraping.Result{Source: s.Name()}
			return
		}

		if response.Data.Total > 0 {
			for _, item := range response.Data.Arr {
				if item.Domain != "" {
					results <- subscraping.Result{Source: s.Name(), Value: item.Domain}
				}
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "hunter"
}
