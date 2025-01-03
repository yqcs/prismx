package healthcheck

import (
	"context"
	"net"
	"strings"
)

type DnsResolveInfo struct {
	Host        string
	Resolver    string
	Successful  bool
	IPAddresses []net.IPAddr
	Error       error
}

func DnsResolve(host string, resolver string) DnsResolveInfo {
	ipAddresses, err := getIPAddresses(host, resolver)

	return DnsResolveInfo{
		Host:        host,
		Resolver:    resolver,
		Successful:  err == nil,
		IPAddresses: ipAddresses,
		Error:       err,
	}
}

func getIPAddresses(name, dnsServer string) ([]net.IPAddr, error) {
	if !strings.Contains(dnsServer, ":") {
		dnsServer = dnsServer + ":53"
	}

	resolver := net.Resolver{
		PreferGo: true, Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, dnsServer)
		}}

	resolvedIPs, err := resolver.LookupIPAddr(context.Background(), name)
	if err != nil {
		return nil, err
	}

	return resolvedIPs, nil
}
