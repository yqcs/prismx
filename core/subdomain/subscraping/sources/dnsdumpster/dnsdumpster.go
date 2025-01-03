// Package dnsdumpster logic
package dnsdumpster

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"prismx_cli/core/subdomain/subscraping"
	"regexp"
	"strings"
)

// CSRFSubMatchLength CSRF regex submatch length
const CSRFSubMatchLength = 2

var re = regexp.MustCompile("<input type=\"hidden\" name=\"csrfmiddlewaretoken\" value=\"(.*)\">")

// getCSRFToken gets the CSRF Token from the page
func getCSRFToken(page string) string {
	if subs := re.FindStringSubmatch(page); len(subs) == CSRFSubMatchLength {
		return strings.TrimSpace(subs[1])
	}
	return ""
}

// postForm posts a form for a domain and returns the response
func postForm(ctx context.Context, session *subscraping.Session, token, domain string) (string, error) {
	params := url.Values{
		"csrfmiddlewaretoken": {token},
		"targetip":            {domain},
		"user":                {"free"},
	}

	resp, err := session.HTTPRequest(
		ctx,
		"POST",
		"https://dnsdumpster.com/",
		fmt.Sprintf("csrftoken=%s; Domain=dnsdumpster.com", token),
		map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Referer":      "https://dnsdumpster.com",
			"X-CSRF-Token": token,
		},
		strings.NewReader(params.Encode()))

	if err != nil {
		return "", err
	}

	// Now, grab the entire page
	in, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	return string(in), err
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)

	go func() {
		defer close(results)

		resp, err := session.SimpleGet(ctx, "https://dnsdumpster.com/")

		if err != nil {
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		csrfToken := getCSRFToken(string(body))
		data, err := postForm(ctx, session, csrfToken, domain)
		if err != nil {
			return
		}

		for _, subdomain := range session.Extractor.FindAllString(data, -1) {
			results <- subscraping.Result{Source: s.Name(), Value: subdomain}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "dnsdumpster"
}
