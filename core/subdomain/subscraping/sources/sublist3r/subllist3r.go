// Package sublist3r logic
package sublist3r

import (
	"context"
	"encoding/json"
	"fmt"
	"prismx_cli/core/subdomain/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://api.sublist3r.com/search.php?domain=%s", domain))
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var subdomains []string
		err = json.NewDecoder(resp.Body).Decode(&subdomains)
		if err != nil {
			return
		}

		for _, subdomain := range subdomains {
			results <- subscraping.Result{Source: s.Name(), Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "sublist3r"
}
