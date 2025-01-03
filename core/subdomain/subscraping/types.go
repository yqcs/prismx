package subscraping

import (
	"context"
	"regexp"
	"time"

	"go.uber.org/ratelimit"
)

type Source interface {
	Run(context.Context, string, *Session) <-chan Result
	Name() string
}

// Session is the option passed to the source, an option is created
type Session struct {
	Timeout     time.Duration
	Extractor   *regexp.Regexp
	Keys        *Keys
	RateLimiter ratelimit.Limiter
}

// Keys contains the current API Keys we have in store
type Keys struct {
	Shodan          string
	ThreatBook      string
	Virustotal      string
	ZoomEyeUserName string
	ZoomEyePass     string
	FofaUsername    string
	FofaSecret      string
	HunterUserName  string
	HunterKey       string
	FullHunt        string
}

var AppKey Keys

type Result struct {
	Source string
	Value  string
}
