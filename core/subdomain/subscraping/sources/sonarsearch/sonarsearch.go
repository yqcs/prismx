// Package sonarsearch logic
package sonarsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"prismx_cli/core/subdomain/subscraping"
	"strconv"
)

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)
	go func() {
		defer close(results)

		getURL := fmt.Sprintf("https://sonar.omnisint.io/subdomains/%s?page=", domain)
		page := 0
		var subdomains []string
		for {
			resp, err := session.SimpleGet(ctx, getURL+strconv.Itoa(page))

			if err != nil {
				results <- subscraping.Result{Source: s.Name()}
				return
			}
			defer resp.Body.Close()
			if err = json.NewDecoder(resp.Body).Decode(&subdomains); err != nil {
				results <- subscraping.Result{Source: s.Name()}
				return
			}

			if len(subdomains) == 0 {
				return
			}

			for _, subdomain := range subdomains {
				results <- subscraping.Result{Source: s.Name(), Value: subdomain}
			}
			page++
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "sonarsearch"
}
