// Package commoncrawl logic
package commoncrawl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"prismx_cli/core/subdomain/subscraping"
	"strings"
)

const indexURL = "https://index.commoncrawl.org/collinfo.json"

type indexResponse struct {
	ID     string `json:"id"`
	APIURL string `json:"cdx-api"`
}

// Source is the passive scraping agent
type Source struct{}

var years = [...]string{"2020", "2019", "2018", "2017"}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, indexURL)

		if err != nil {
			return
		}
		defer resp.Body.Close()
		var indexes []indexResponse
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &indexes); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		searchIndexes := make(map[string]string)
		for _, year := range years {
			for _, index := range indexes {
				if strings.Contains(index.ID, year) {
					if _, ok := searchIndexes[year]; !ok {
						searchIndexes[year] = index.APIURL
						break
					}
				}
			}
		}

		for _, apiURL := range searchIndexes {
			further := s.getSubdomains(ctx, apiURL, domain, session, results)
			if !further {
				break
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "commoncrawl"
}

func (s *Source) getSubdomains(ctx context.Context, searchURL, domain string, session *subscraping.Session, results chan subscraping.Result) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			var headers = map[string]string{"Host": "index.commoncrawl.org"}
			resp, err := session.Get(ctx, fmt.Sprintf("%s?url=*.%s", searchURL, domain), "", headers)
			if err != nil {
				results <- subscraping.Result{Source: s.Name()}

				return false
			}

			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				line, _ = url.QueryUnescape(line)
				subdomain := session.Extractor.FindString(line)
				if subdomain != "" {
					// fix for triple encoded URL
					subdomain = strings.ToLower(subdomain)
					subdomain = strings.TrimPrefix(subdomain, "25")
					subdomain = strings.TrimPrefix(subdomain, "2f")

					results <- subscraping.Result{Source: s.Name(), Value: subdomain}
				}
			}
			resp.Body.Close()
			return true
		}
	}
}
