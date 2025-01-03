//go:build linux
// +build linux

package updateutils

import (
	"encoding/base64"
	"runtime"
	"strings"

	"github.com/zcalusic/sysinfo"
)

// Get OS Vendor returns the linux distribution vendor
// if not linux then returns runtime.GOOS
func GetOSVendor() string {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	return si.OS.Vendor
}

// returns platform metadata
func getPlatformMetadata() string {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	tmp := strings.ReplaceAll(si.Board.Vendor, " ", "_") + "|" + strings.ReplaceAll(si.Board.Name, " ", "_")
	if tmp == "|" {
		// instead of just empty string return os for more context
		tmp = runtime.GOOS + "|" + runtime.GOARCH
	}
	return strings.TrimSuffix(base64.StdEncoding.EncodeToString([]byte(tmp)), "==")
}
