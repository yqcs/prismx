package passive

import (
	"context"
	"prismx_cli/core/subdomain/subscraping"
	"strings"
	"sync"
	"time"
)

// EnumerateSubdomains enumerates all the subdomains for a given domain
func (a *Agent) EnumerateSubdomains(domain string, keys *subscraping.Keys, maxEnumTime time.Duration) ([]subscraping.Result, error) {

	session, err := subscraping.NewSession(domain, keys, maxEnumTime)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxEnumTime)
	wg := &sync.WaitGroup{}
	var results []subscraping.Result

	// 来源目标。
	for source, runner := range a.sources {
		wg.Add(1)
		go func(source string, runner subscraping.Source) {
			for resp := range runner.Run(ctx, domain, session) {
				resp.Value = strings.ToLower(resp.Value)
				results = append(results, resp)
			}
			wg.Done()
		}(source, runner)
	}
	wg.Wait()
	cancel()
	return results, nil
}
