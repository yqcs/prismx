// Package crtsh logic
package crtsh

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"prismx_cli/core/subdomain/subscraping"
	"prismx_cli/utils/arr"
	"strings"
)

type subdomain struct {
	ID        int    `json:"id"`
	NameValue string `json:"name_value"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)
		resp, err := session.SimpleGet(ctx, fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain))
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		defer resp.Body.Close()
		var subdomains []subdomain
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &subdomains); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		var doList []string
		for _, item := range subdomains {
			doList = append(doList, strings.Split(item.NameValue, "\n")...)
		}
		doList = arr.SliceRemoveDuplicates(doList)
		for _, item := range doList {
			results <- subscraping.Result{Source: s.Name(), Value: item}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "crtsh"
}
