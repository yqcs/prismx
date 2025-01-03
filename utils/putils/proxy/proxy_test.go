//go:build proxy

package proxyutils

// package tests will be executed only with (running proxy is necessary):
// go test -tags proxy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const burpURL = "http://127.0.0.1:8080"

// a local instance of burp community is necessary
func TestIsBurp(t *testing.T) {
	ok, err := IsBurp(burpURL)
	require.Nil(t, err)
	require.True(t, ok)
}

// a valid proxy is necessary
func TestValidateOne(t *testing.T) {
	proxyURL, err := ValidateOne(burpURL)
	require.Nil(t, err)
	require.Equal(t, burpURL, proxyURL)
}
