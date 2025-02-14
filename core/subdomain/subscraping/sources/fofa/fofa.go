// Package fofa logic
package fofa

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
	"strings"
)

type fofaResponse struct {
	Error   bool     `json:"error"`
	ErrMsg  string   `json:"errmsg"`
	Size    int      `json:"size"`
	Results []string `json:"results"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.FofaUsername == "" || session.Keys.FofaSecret == "" {
			return
		}

		// fofa api doc https://fofa.info/static_pages/api_help
		qbase64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("domain=\"%s\"", domain)))
		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://fofa.info/api/v1/search/all?full=true&fields=host&page=1&size=10000&email=%s&key=%s&qbase64=%s", session.Keys.FofaUsername, session.Keys.FofaSecret, qbase64))
		if err != nil && resp == nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		defer resp.Body.Close()
		var response fofaResponse
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &response); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if response.Error {
			results <- subscraping.Result{Source: s.Name()}
			return
		}

		if response.Size > 0 {
			for _, subdomain := range response.Results {
				if strings.HasPrefix(strings.ToLower(subdomain), "http://") || strings.HasPrefix(strings.ToLower(subdomain), "https://") {
					subdomain = subdomain[strings.Index(subdomain, "//")+2:]
				}
				results <- subscraping.Result{Source: s.Name(), Value: subdomain}
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "fofa"
}
