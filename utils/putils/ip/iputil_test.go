package iputil

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"prismx_cli/utils/putils/consts"
	osutils "prismx_cli/utils/putils/os"
)

func TestTryExtendIP(t *testing.T) {
	if osutils.IsLinux() {
		return
	}
	if osutils.IsWindows() {
		i, err := TryExtendIP("localhost")
		require.Nil(t, i)
		require.ErrorIs(t, err, consts.ErrNotSupported)
		return
	}

	type extendIPTestCase struct {
		input         string
		expectedIP    net.IP
		expectedError bool
	}

	testCases := []extendIPTestCase{
		{
			input:         "127.0.0.1:80",
			expectedIP:    net.ParseIP("127.0.0.1"),
			expectedError: false,
		},
		{
			input:         "localhost:1",
			expectedIP:    net.ParseIP("127.0.0.1"),
			expectedError: false,
		},
		{
			input:         "invalid-ip:80",
			expectedIP:    nil,
			expectedError: true,
		},
		{
			input:         "35.1",
			expectedIP:    net.ParseIP("35.0.0.1"),
			expectedError: false,
		},
		{
			input:         "35.1.124",
			expectedIP:    net.ParseIP("35.1.0.124"),
			expectedError: false,
		},
		{
			input:         "192.168.1",
			expectedIP:    net.ParseIP("192.168.0.1"),
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		ip, err := TryExtendIP(tc.input)
		require.Equal(t, tc.expectedError, err != nil, "input: %v, error: %v", tc.input, err)
		require.True(t, ip.Equal(tc.expectedIP), "Expected IP: %v, got: %v", tc.expectedIP, ip)
	}
}

func TestCanExtend(t *testing.T) {
	if osutils.IsWindows() || osutils.IsLinux() {
		return
	}

	tests := map[string]bool{
		"35.1":          true,
		"0":             true,
		"1":             true,
		"199":           true,
		"zippo":         false,
		"-1":            false,
		"1 1":           false,
		"localhost":     true,
		"192.168.1.1":   true,
		"192.168.1.547": false,
		"192.168.1.258": false,
		"192.168.256.1": false,
		"1.1.1.1.1":     false,
	}
	for item, expected := range tests {
		got := CanExtend(item)
		require.Equal(t, expected, got, "Expected: %v, got: %v => %v", expected, got, item)
	}
}

func TestIsIPv6Short(t *testing.T) {
	type test struct {
		Ip           string
		Expected     bool
		MessageError string
	}

	validIpsTest := []test{
		{
			Ip:           "001:0db8::1",
			Expected:     true,
			MessageError: "valid ip not recognized",
		},
		{
			Ip:           "::a00:27ff:fef3:7d56",
			Expected:     true,
			MessageError: "valid ip not recognized",
		},
		{
			Ip:           "fe80::a00:27ff:fef3:7d56",
			Expected:     true,
			MessageError: "valid ip not recognized",
		},
		{
			Ip:           "2607:f0d0:1002:0051:0000:0000:0000:0004",
			Expected:     true,
			MessageError: "valid ip not recognized",
		},
	}

	for _, ip := range validIpsTest {
		require.Equal(t, ip.Expected, IsIPv6(ip.Ip), ip.MessageError, ip.Ip)
	}
}

func TestIsInternalIPv4(t *testing.T) {
	// Test this ipv4
	require.False(t, IsInternal("153.12.14.1"), "internal ipv4 address recognized as not valid")
	require.True(t, IsInternal("172.16.0.0"), "internal ipv4 address recognized as valid")
	// Test with ipv6
	require.False(t, IsInternal("684D:1111:222:3333:4444:5555:6:77"), "internal ipv4 address recognized as not valid")
	require.True(t, IsInternal("fc00:7e5b:cfa9::"), "internal ipv4 address recognized as valid")
}

func TestIsPort(t *testing.T) {
	require.False(t, IsPort("0"), "invalid port 0")
	require.False(t, IsPort("-1"), "negative port")
	require.True(t, IsPort("1"), "valid port not recognized")
	require.True(t, IsPort("65535"), "valid port not recognized")
	require.False(t, IsPort("65536"), "valid port not recognized")
	require.False(t, IsPort("0xff"), "hex port considered valid")
	require.False(t, IsPort("12.12"), "float recognized as valid")
}

func TestIsIPv4(t *testing.T) {
	require.True(t, IsIPv4("127.0.0.1"), "valid ipv4 address not recognized")
	require.False(t, IsIPv4("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), "ipv6 address recognized as valid")
}

func TestIsIPv6(t *testing.T) {
	require.False(t, IsIPv6("127.0.0.1"), "ipv4 address recognized as valid")
	require.True(t, IsIPv6("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), "valid ipv6 address not recognized")
	require.True(t, IsIPv6("::a00:27ff:fef3:7d56"), "valid ipv6 address not recognized")
	require.True(t, IsIPv6("001:0db8::1"), "valid ipv6 address not recognized")
	require.True(t, IsIPv6("::1"), "valid ipv6 address not recognized")
}

func TestIsCIDR(t *testing.T) {
	require.False(t, IsCIDR("127.0.0.1"), "ipv4 address recognized as cidr")
	require.True(t, IsCIDR("127.0.0.0/24"), "valid cidr not recognized")
	require.True(t, IsCIDR("127.0.0.0/1"), "valid cidr not recognized")
	require.True(t, IsCIDR("127.0.0.0/32"), "valid cidr not recognized")
	require.False(t, IsCIDR("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), "ipv6 address recognized as cidr")
}

func TestIsCidrWithExpansion(t *testing.T) {
	require.True(t, IsCidrWithExpansion("127.0.0.1-32"), "valid cidr /32 not recognized")
	require.False(t, IsCidrWithExpansion("127.0.0.0-55"), "invalid cidr /55")
}

func TestToCidr(t *testing.T) {
	tests := map[string]bool{
		"127.0.0.0/24": true,
		"127.0.0.1":    true,
		"aaa":          false,
	}
	for item, ok := range tests {
		tocidr := ToCidr(item)
		if ok {
			require.NotNil(t, tocidr, "valid cidr not recognized")
		} else {
			require.Nil(t, tocidr, "invalid cidr")
		}
	}
}

func TestAsIPV4IpNet(t *testing.T) {
	tests := map[string]bool{
		"127.0.0.0/24": true,
		"127.0.0.1":    true,
		"aaa":          false,
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334": false,
	}
	for item, ok := range tests {
		tocidr := AsIPV4IpNet(item)
		if ok {
			require.NotNil(t, tocidr, "valid cidr not recognized")
		} else {
			require.Nil(t, tocidr, "invalid cidr")
		}
	}
}

func TestAsIPV6IpNet(t *testing.T) {
	tests := map[string]bool{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334": true,
		"2002::1234:abcd:ffff:c0a8:101/64":        true,
	}
	for item, ok := range tests {
		tocidr := AsIPV6CIDR(item)
		if ok {
			require.NotNil(t, tocidr, "valid cidr not recognized")
		} else {
			require.Nil(t, tocidr, "invalid cidr")
		}
	}
}

func TestWhatsMyIP(t *testing.T) {
	// we can't compare the ip with local interfaces as it might be the external gateway one
	// so we just verify we can contact the api endpoint
	_, err := WhatsMyIP()
	require.Nil(t, err, "couldn't retrieve ip")
}

func TestToFQDN(t *testing.T) {
	// we can't compare the ip with local interfaces as it might be the external gateway one
	// so we just verify we can contact the api endpoint
	fqdns, err := ToFQDN("1.1.1.1")
	require.Nil(t, err, "couldn't retrieve ip")
	require.NotNil(t, fqdns)
}
