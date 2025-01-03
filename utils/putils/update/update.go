package updateutils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"prismx_cli/utils/machineid"
	"runtime"
	"time"
)

const (
	Organization        = "projectdiscovery"
	UpdateCheckEndpoint = "https://api.pdtm.sh/api/v1/tools/%v"
)

var (
	// By default when tool is updated release notes of latest version are printed
	HideReleaseNotes      = false
	HideProgressBar       = false
	VersionCheckTimeout   = time.Duration(5) * time.Second
	DownloadUpdateTimeout = time.Duration(30) * time.Second
	// Note: DefaultHttpClient is only used in GetToolVersionCallback
	DefaultHttpClient *http.Client
)

// GetpdtmParams returns encoded query parameters sent to update check endpoint
func GetpdtmParams(version string) string {
	params := &url.Values{}
	os := runtime.GOOS
	if runtime.GOOS == "linux" {
		// be more specific
		os = GetOSVendor()
	}
	params.Add("os", os)
	params.Add("arch", runtime.GOARCH)
	params.Add("go_version", runtime.Version())
	params.Add("v", version)
	params.Add("machine_id", GetMachineID())
	params.Add("utm_source", getUtmSource())
	return params.Encode()
}

// GetMachineID return a unique identifier that is unique to the machine
// it is a sha256 hashed value with pdtm as salt
func GetMachineID() string {
	machineId, err := machineid.ProtectedID("pdtm")
	if err != nil {
		return "unknown"
	}
	return machineId
}

func init() {
	DefaultHttpClient = &http.Client{
		Timeout: VersionCheckTimeout,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
