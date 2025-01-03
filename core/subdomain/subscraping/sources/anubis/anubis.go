// Package anubis logic
package anubis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://jonlu.ca/anubis/subdomains/%s", domain))
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		defer resp.Body.Close()
		var subdomains []string
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &subdomains); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		for _, record := range subdomains {
			results <- subscraping.Result{Source: s.Name(), Value: record}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "anubis"
}
