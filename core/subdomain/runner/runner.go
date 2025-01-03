package runner

import (
	"prismx_cli/core/subdomain/passive"
	"prismx_cli/core/subdomain/subscraping"
	"time"
)

type Runner struct {
	Target  string
	Timeout time.Duration
}

func RunEnumeration(run Runner) ([]subscraping.Result, error) {
	agent := passive.New(passive.DefaultAllSources, []string{})
	passiveResults, err := agent.EnumerateSubdomains(run.Target, &subscraping.AppKey, run.Timeout*3)
	return passiveResults, err
}
