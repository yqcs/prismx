package healthcheck

import (
	"path/filepath"

	fileutil "prismx_cli/utils/putils/file"
	folderutil "prismx_cli/utils/putils/folder"
)

var (
	DefaultPathsToCheckPermission   = []string{filepath.Join(folderutil.HomeDirOrDefault(""), ".config", fileutil.ExecutableName())}
	DefaultHostsToCheckConnectivity = []string{"scanme.sh"}
	DefaultResolver                 = "1.1.1.1:53"
)

type HealthCheckInfo struct {
	EnvironmentInfo EnvironmentInfo
	PathPermissions []PathPermission
	DnsResolveInfos []DnsResolveInfo
}

type Options struct {
	Paths    []string
	Hosts    []string
	Resolver string
}

var DefaultOptions = Options{
	Paths:    DefaultPathsToCheckPermission,
	Hosts:    DefaultHostsToCheckConnectivity,
	Resolver: DefaultResolver,
}

func Do(programVersion string, options *Options) (healthCheckInfo HealthCheckInfo) {
	if options == nil {
		options = &DefaultOptions
	}
	healthCheckInfo.EnvironmentInfo = CollectEnvironmentInfo(programVersion)
	for _, path := range options.Paths {
		healthCheckInfo.PathPermissions = append(healthCheckInfo.PathPermissions, CheckPathPermission(path))
	}
	for _, host := range options.Hosts {
		healthCheckInfo.DnsResolveInfos = append(healthCheckInfo.DnsResolveInfos, DnsResolve(host, options.Resolver))
	}
	return
}
