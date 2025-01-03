package fullhunt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"prismx_cli/core/subdomain/subscraping"
)

// fullhunt response
type fullHuntResponse struct {
	Hosts   []string `json:"hosts"`
	Message string   `json:"message"`
	Status  int      `json:"status"`
}

// Source is the passive scraping agent
type Source struct{}

func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.Get(ctx, fmt.Sprintf("https://fullhunt.io/api/v1/domain/%s/subdomains", domain), "", map[string]string{"X-API-KEY": session.Keys.FullHunt})

		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		defer resp.Body.Close()
		var response fullHuntResponse
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		if err = json.Unmarshal(all, &response); err != nil {
			results <- subscraping.Result{Source: s.Name()}
			return
		}
		for _, record := range response.Hosts {
			results <- subscraping.Result{Source: s.Name(), Value: record}
		}
	}()
	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "fullhunt"
}
