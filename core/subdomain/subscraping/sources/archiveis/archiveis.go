// Package archiveis is a Archiveis Scraping Engine in Golang
package archiveis

import (
	"context"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
)

// Source is the passive scraping agent
type Source struct{}

// Name returns the name of the source
func (s *Source) Name() string {
	return "archiveis"
}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://archive.is/*.%s", domain))
		if err != nil {
			results <- subscraping.Result{Source: "archiveis"}
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: "archiveis"}
			return
		}

		src := string(body)

		for _, subdomain := range session.Extractor.FindAllString(src, -1) {
			results <- subscraping.Result{Source: "archiveis", Value: subdomain}
		}
	}()

	return results
}
