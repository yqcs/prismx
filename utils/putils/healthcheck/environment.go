package healthcheck

import (
	"os"
	"runtime"

	"github.com/projectdiscovery/fdmax"
	iputil "prismx_cli/utils/putils/ip"
	permissionutil "prismx_cli/utils/putils/permission"
	router "prismx_cli/utils/putils/routing"
)

type EnvironmentInfo struct {
	ExternalIPv4   string
	Admin          bool
	Arch           string
	Compiler       string
	GoVersion      string
	OSName         string
	ProgramVersion string
	OutboundIPv4   string
	OutboundIPv6   string
	Ulimit         Ulimit
	PathEnvVar     string
	Error          error
}

type Ulimit struct {
	Current uint64
	Max     uint64
}

func CollectEnvironmentInfo(appVersion string) EnvironmentInfo {
	externalIPv4, _ := iputil.WhatsMyIP()
	outboundIPv4, outboundIPv6, _ := router.GetOutboundIPs()

	ulimit := Ulimit{}
	limit, err := fdmax.Get()
	if err == nil {
		ulimit.Current = limit.Current
		ulimit.Max = limit.Max
	}

	return EnvironmentInfo{
		ExternalIPv4:   externalIPv4,
		Admin:          permissionutil.IsRoot,
		Arch:           runtime.GOARCH,
		Compiler:       runtime.Compiler,
		GoVersion:      runtime.Version(),
		OSName:         runtime.GOOS,
		ProgramVersion: appVersion,
		OutboundIPv4:   outboundIPv4.String(),
		OutboundIPv6:   outboundIPv6.String(),
		Ulimit:         ulimit,
		PathEnvVar:     os.Getenv("PATH"),
	}
}
