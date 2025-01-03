// Package virustotal logic
package virustotal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"prismx_cli/core/subdomain/subscraping"
)

type response struct {
	Subdomains []string `json:"subdomains"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.Virustotal == "" {
			return
		}

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", session.Keys.Virustotal, domain))

		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
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

		for _, subdomain := range data.Subdomains {
			results <- subscraping.Result{Source: s.Name(), Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "virustotal"
}
