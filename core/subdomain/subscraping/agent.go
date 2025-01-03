package subscraping

import (
	"context"
	"io"
	"net/http"
	"prismx_cli/utils/netUtils"
	"regexp"
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

var subdomainExtractorMutex = &sync.Mutex{}

// NewSession creates a new session object for a domain
func NewSession(domain string, keys *Keys, timeout time.Duration) (*Session, error) {
	session := &Session{
		Keys:    keys,
		Timeout: timeout,
	}
	session.RateLimiter = ratelimit.NewUnlimited()
	subdomainExtractorMutex.Lock()
	extractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
	subdomainExtractorMutex.Unlock()

	if err != nil {
		return nil, err
	}
	session.Extractor = extractor

	return session, err
}

// Get makes a GET request to a URL with extended parameters
func (s *Session) Get(ctx context.Context, getURL, cookies string, headers map[string]string) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodGet, getURL, cookies, headers, nil)
}

// SimpleGet makes a simple GET request to a URL
func (s *Session) SimpleGet(ctx context.Context, getURL string) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodGet, getURL, "", map[string]string{}, nil)
}

// Post makes a POST request to a URL with extended parameters
func (s *Session) Post(ctx context.Context, postURL, cookies string, headers map[string]string, body io.Reader) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodPost, postURL, cookies, headers, body)
}

// SimplePost makes a simple POST request to a URL
func (s *Session) SimplePost(ctx context.Context, postURL, contentType string, body io.Reader) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodPost, postURL, "", map[string]string{"Content-Type": contentType}, body)
}

// HTTPRequest makes any HTTP request to a URL with extended parameters
func (s *Session) HTTPRequest(ctx context.Context, method, requestURL, cookies string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Connection", "close")

	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	s.RateLimiter.Take()
	do, err := netUtils.SendHttp(req, s.Timeout, true)
	if err != nil {
		return nil, err
	}
	return do.Other, err
}
