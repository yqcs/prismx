// Package threatminer logic
package threatminer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"prismx_cli/core/subdomain/subscraping"
)

type response struct {
	StatusCode    string   `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Results       []string `json:"results"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://api.threatminer.org/v2/domain.php?q=%s&rt=5", domain))

		if err != nil {
			return
		}
		defer resp.Body.Close()
		var data response
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &data); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		for _, subdomain := range data.Results {
			results <- subscraping.Result{Source: s.Name(), Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "threatminer"
}
