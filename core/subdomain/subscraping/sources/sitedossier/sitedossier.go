// Package sitedossier logic
package sitedossier

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"prismx_cli/core/subdomain/subscraping"
	"regexp"
	"time"
)

// SleepRandIntn is the integer value to get the pseudo-random number
// to sleep before find the next match
const SleepRandIntn = 5

var reNext = regexp.MustCompile(`<a href="([A-Za-z0-9/.]+)"><b>`)

type agent struct {
	results chan subscraping.Result
	session *subscraping.Session
}

func (a *agent) enumerate(ctx context.Context, baseURL string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	resp, err := a.session.SimpleGet(ctx, baseURL)
	if err != nil {
		a.results <- subscraping.Result{Source: "sitedossier"}
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.results <- subscraping.Result{Source: "sitedossier"}
		return
	}

	src := string(body)
	for _, match := range a.session.Extractor.FindAllString(src, -1) {
		a.results <- subscraping.Result{Source: "sitedossier", Value: match}
	}

	match1 := reNext.FindStringSubmatch(src)
	time.Sleep(time.Duration((3 + rand.Intn(SleepRandIntn))) * time.Second)

	if len(match1) > 0 {
		a.enumerate(ctx, "http://www.sitedossier.com"+match1[1])
	}
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	a := agent{
		session: session,
		results: results,
	}

	go func() {
		a.enumerate(ctx, fmt.Sprintf("http://www.sitedossier.com/parentdomain/%s", domain))
		close(a.results)
	}()

	return a.results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "sitedossier"
}
