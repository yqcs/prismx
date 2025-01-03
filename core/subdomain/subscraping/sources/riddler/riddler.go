// Package riddler logic
package riddler

import (
	"bufio"
	"context"
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

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://riddler.io/search?q=pld:%s&view_type=data_table", domain))

		if err != nil {
			return
		}
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			subdomain := session.Extractor.FindString(line)
			if subdomain != "" {
				results <- subscraping.Result{Source: s.Name(), Value: subdomain}
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "riddler"
}
