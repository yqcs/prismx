package urlutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// full url
	U, err := Parse("http://127.0.0.1/a")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "http", U.Scheme, "different scheme")
	require.Equal(t, "127.0.0.1", U.Hostname(), "different host")
	require.Equal(t, "/a", U.Path, "different request uri")

	// full url with port
	U, err = Parse("http://127.0.0.1:1000/a")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "http", U.Scheme, "different scheme")
	require.Equal(t, "127.0.0.1", U.Hostname(), "different host")
	require.Equal(t, "1000", U.Port(), "different host")
	require.Equal(t, "/a", U.Path, "different request uri")

	// partial url without port
	U, err = Parse("a.b.c.d")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "a.b.c.d", U.Hostname(), "different host")
	require.Equal(t, "", U.Path, "different request uri")

	// partial url with protocol and no port
	U, err = Parse("https://a.b.c.d")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "https", U.Scheme, "different scheme")
	require.Equal(t, "a.b.c.d", U.Hostname(), "different host")
	require.Equal(t, "", U.Path, "different request uri")

	// replacing port
	U, err = Parse("https://a.b.c.d")
	require.Nil(t, err, "could not parse url")
	U.UpdatePort("15000")
	require.Equal(t, "https://a.b.c.d:15000", U.String(), "port not replaced")

	// replacing port
	U, err = Parse("https://a.b.c.d//d")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "https://a.b.c.d//d", U.String(), "unexpected url")

	// fragmented url
	U, err = Parse("http://127.0.0.1/#a")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "http", U.Scheme, "different scheme")
	require.Equal(t, "127.0.0.1", U.Hostname(), "different host")
	require.Equal(t, "a", U.Fragment, "different fragment")
	require.Equal(t, "http://127.0.0.1/#a", U.String(), "different full url")

	// websocket
	U, err = Parse("wss://127.0.0.1")
	require.Nil(t, err, "could not parse url")
	require.Equal(t, "wss", U.Scheme, "different scheme")
	require.Equal(t, "127.0.0.1", U.Hostname(), "different host")
	require.Equal(t, "wss://127.0.0.1", U.String(), "different full url")
}

func TestClone(t *testing.T) {
	U, err := Parse("https://scanme.sh/some/path?with=param#fragment")
	require.Nil(t, err)
	cloned := U.Clone()
	require.Equal(t, U, cloned)

	U, err = Parse("https://secret:secret@scanme.sh/some/path?with=param#fragment")
	require.Nil(t, err)
	cloned = U.Clone()
	require.Equal(t, U, cloned)
}

func TestPortUpdate(t *testing.T) {
	expected := "http://localhost:8000/test"
	urlx, err := Parse("http://localhost:53/test")
	require.Nil(t, err)
	urlx.UpdatePort("8000")
	require.Equalf(t, urlx.String(), expected, "expected %v but got %v", expected, urlx.String())
}

func TestUpdateRelPath(t *testing.T) {
	// updates existing relative path with new one
	exURL := "https://scanme.sh/somepath/abc?key=true"
	urlx, err := Parse(exURL)
	require.Nil(t, err)
	err = urlx.UpdateRelPath("/newpath/?with=params", true)
	require.Nil(t, err)
	require.Equalf(t, urlx.Path, "/newpath/", "failed to update relative path")
}

func TestInvalidURLs(t *testing.T) {
	testcases := []string{
		"https://scanme.sh/%invalid/%0D%0A",
		"https://scanme.sh/%invalid2/and/path",
		"https://scanme.sh",
		"https://scanme.sh/%invalid?with=param",
		"https://127.0.0.1:52272/%invalid",
		"http.s3.amazonaws.com",
		"https.s3.amazonaws.com",
		"scanme.sh/xyz/invalid",
		"scanme.sh/xyz/%u2s/%invalid",
	}
	for _, v := range testcases {
		urlx, err := ParseAbsoluteURL(v, true)
		require.Nilf(t, err, "got error for url %v", v)
		require.Equal(t, v, urlx.String())
	}
}

func TestParseRelativePath(t *testing.T) {
	testcases := []struct {
		inputURL     string
		unsafe       bool
		expectedPath string
	}{
		{"//CFIDE/wizards/common/utils.cfc", false, "//CFIDE/wizards/common/utils.cfc"},
	}

	for _, v := range testcases {
		urlx, err := ParseRelativePath(v.inputURL, v.unsafe)
		require.Nilf(t, err, "got error for url %v", v.inputURL)
		require.Equal(t, v.expectedPath, urlx.GetRelativePath())
	}
}

func TestParseFragmentRelativePath(t *testing.T) {
	testcases := []struct {
		inputURL         string
		unsafe           bool
		expectedFragment string
	}{
		{"/?param=value#highlight", false, "highlight"},
		{"/somepath/#highlight", true, "highlight"},
		{"/somepath/?with=param#highlight", false, "highlight"},
	}

	for _, v := range testcases {
		urlx, err := ParseURL(v.inputURL, v.unsafe)
		require.Nilf(t, err, "got error for url %v", v.inputURL)
		require.Equal(t, v.expectedFragment, urlx.Fragment)
	}
}

func TestParseInvalidUnsafe(t *testing.T) {
	testcases := []string{
		"https://127.0.0.1/%25",
		"https://127.0.0.1/%25/aaaa",
		"https://127.0.0.1/%25/bb/%45?a=1",
	}

	for _, input := range testcases {
		u, err := ParseURL(input, true)
		require.Nilf(t, err, "got error for url %v", input)
		require.Equal(t, input, u.URL.String())
	}
}

func TestParseParam(t *testing.T) {
	testcases := []struct {
		inputURL      string
		unsafe        bool
		expectedParam string
	}{
		{"https://scanme.sh/?param=value#highlight", false, "param=value"},
	}

	for _, v := range testcases {
		urlx, err := ParseURL(v.inputURL, v.unsafe)
		require.Nilf(t, err, "got error for url %v", v.inputURL)
		require.Equal(t, v.expectedParam, urlx.Params.Encode())
	}
}

func TestURLFragment(t *testing.T) {
	testcases := []struct {
		inputURL         string
		unsafe           bool
		expectedFragment string
	}{
		{"https://scanme.sh/admin?param=value#highlight", false, "highlight"},
		{"https://scanme.sh/#highlight", true, "highlight"},
	}

	for _, v := range testcases {
		urlx, err := ParseURL(v.inputURL, v.unsafe)
		require.Nilf(t, err, "got error for url %v", v.inputURL)
		require.Equalf(t, v.expectedFragment, urlx.Fragment, "got error for url %v", v.inputURL)
	}
}

func TestUnicodeEscapeWithUnsafe(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{"https://scanme.sh/%u002e%u002e/%u002e%u002e/1.txt.it", "https://scanme.sh/%u002e%u002e/%u002e%u002e/1.txt.it"},
	}

	for _, v := range testcases {
		urlx, err := ParseAbsoluteURL(v.input, true)
		require.Nilf(t, err, "got error for url %v", v.input)
		require.Equal(t, v.expected, urlx.String())
	}
}

func TestInvalidScheme(t *testing.T) {
	testcases := []struct {
		input       string
		expectedErr bool
	}{
		{"//:foo", true},
		{"://foo", true},
	}
	for _, v := range testcases {
		urlx, err := ParseAbsoluteURL(v.input, true)
		if v.expectedErr {
			require.NotNil(t, err)
			require.Nil(t, urlx)
		} else {
			require.Nil(t, err)
			require.NotNil(t, urlx)
		}
	}
}
