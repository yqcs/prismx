//go:build !linux
// +build !linux

package updateutils

import (
	"encoding/base64"
	"runtime"
	"strings"
)

// Get OS Vendor returns the linux distribution vendor
// if not linux then returns runtime.GOOS
func GetOSVendor() string {
	return runtime.GOOS
}

// returns platform metadata
func getPlatformMetadata() string {
	tmp := runtime.GOOS + "|" + runtime.GOARCH
	return strings.TrimSuffix(base64.StdEncoding.EncodeToString([]byte(tmp)), "==")
}
