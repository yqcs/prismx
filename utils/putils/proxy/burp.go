package proxyutils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// IsBurp checks if the target proxy URL is burp suite
func IsBurp(proxyURL string) (bool, error) {
	return getURLWithHTTPProxy("http://burpsuite/", proxyURL, func(resp *http.Response) (bool, error) {
		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("unexpected status code (200 wanted): %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		defer resp.Body.Close()

		return bytes.Contains(body, []byte("Burp Suite")), nil
	})
}

// ValidateOne returns the first valid proxy from a list of proxies by setting up a test connection with scanme.sh
func ValidateOne(proxies ...string) (string, error) {
	for _, proxy := range proxies {
		ok, err := getURLWithHTTPProxy("https://scanme.sh", proxy, func(resp *http.Response) (bool, error) {
			if resp.StatusCode != http.StatusOK {
				return false, fmt.Errorf("unexpected status code (200 wanted): %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()

			return len(body) > 0, nil
		})
		if ok {
			return proxy, err
		}
	}

	return "", errors.New("no valid proxy found")
}

func getURLWithHTTPProxy(targetURL, proxyURL string, checkCallback func(resp *http.Response) (bool, error)) (bool, error) {
	URL, err := url.Parse(proxyURL)
	if err != nil {
		return false, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyURL(URL),
		},
	}

	resp, err := httpClient.Get(targetURL)
	if err != nil {
		return false, err
	}

	return checkCallback(resp)
}
