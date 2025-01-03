// Package shodan logic
package shodan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

type dnsdbLookupResponse struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Result     int      `json:"result"`
	Error      string   `json:"error"`
}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		if session.Keys.Shodan == "" {
			return
		}

		searchURL := fmt.Sprintf("https://api.shodan.io/dns/domain/%s?key=%s", domain, session.Keys.Shodan)
		resp, err := session.SimpleGet(ctx, searchURL)

		if err != nil {
			return
		}
		defer resp.Body.Close()
		var response dnsdbLookupResponse
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &response); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}

		if response.Error != "" {
			results <- subscraping.Result{Source: s.Name()}
			return
		}

		for _, data := range response.Subdomains {
			results <- subscraping.Result{Source: s.Name(), Value: fmt.Sprintf("%s.%s", data, domain)}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "shodan"
}
