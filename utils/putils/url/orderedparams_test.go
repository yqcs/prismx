package urlutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderedParam(t *testing.T) {
	p := NewOrderedParams()
	p.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	p.Add("xss", "<script>alert('XSS')</script>")
	p.Add("xssiwthspace", "<svg id=alert(1) onload=eval(id)>")
	p.Add("jsprotocol", "javascript://alert(1)")
	// Note keys are sorted
	expected := "sqli=1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)&xss=<script>alert('XSS')</script>&xssiwthspace=<svg+id=alert(1)+onload=eval(id)>&jsprotocol=javascript://alert(1)"
	require.Equalf(t, expected, p.Encode(), "failed to encode parameters expected %v but got %v", expected, p.Encode())
}

// TestOrderedParamIntegration preserves order of parameters
// while sending request to server (ref:https://github.com/projectdiscovery/nuclei/issues/3801)
func TestOrderedParamIntegration(t *testing.T) {
	expected := "/?xss=<script>alert('XSS')</script>&sqli=1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)&jsprotocol=javascript://alert(1)&xssiwthspace=<svg+id=alert(1)+onload=eval(id)>"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equalf(t, expected, r.RequestURI, "expected %v but got %v", expected, r.RequestURI)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	p := NewOrderedParams()
	p.Add("xss", "<script>alert('XSS')</script>")
	p.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	p.Add("jsprotocol", "javascript://alert(1)")
	p.Add("xssiwthspace", "<svg id=alert(1) onload=eval(id)>")

	url, err := url.Parse(srv.URL)
	require.Nil(t, err)
	url.RawQuery = p.Encode()
	_, err = http.Get(url.String())
	require.Nil(t, err)
}

func TestGetOrderedParams(t *testing.T) {
	values := url.Values{}
	values.Add("sqli", "1+AND+(SELECT+*+FROM+(SELECT(SLEEP(12)))nQIP)")
	values.Add("xss", "<script>alert('XSS')</script>")
	p := GetParams(values)
	require.NotNilf(t, p, "expected params but got nil")
	require.Equalf(t, p.Get("sqli"), values.Get("sqli"), "malformed or missing value for param sqli expected %v but got %v", values.Get("sqli"), p.Get("sqli"))
	require.Equalf(t, p.Get("xss"), values.Get("xss"), "malformed or missing value for param xss expected %v but got %v", values.Get("xss"), p.Get("xss"))
}

func TestIncludeEquals(t *testing.T) {
	p := NewOrderedParams()
	p.Add("key1", "")
	p.Add("key2", "value2")
	if encoded := p.Encode(); encoded != "key1&key2=value2" {
		t.Errorf("Expected 'key1&key2=value2', got '%s'", encoded)
	}

	p.IncludeEquals = true
	if encoded := p.Encode(); encoded != "key1=&key2=value2" {
		t.Errorf("Expected 'key1=&key2=value2', got '%s'", encoded)
	}
}
