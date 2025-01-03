package parse

import (
	"errors"
	"net"
	"prismx_cli/utils/arr"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var errTips = errors.New(" host parsing error\n" +
	"format: \n" +
	"prismx.io\n" +
	"192.168.1.1\n" +
	"192.168.1.1/8\n" +
	"192.168.1.1/16\n" +
	"192.168.1.1/24\n" +
	"192.168.1.1,192.168.1.2\n" +
	"192.168.1.1-192.168.255.255\n" +
	"192.168.1.1-255")

func ParseIP(ip string, noHost string) (hosts, domainList []string, err error) {
	if ip != "" {
		hosts, domainList, err = ParseIPs(ip)
		if err != nil {
			return nil, nil, err
		}
	}
	if noHost != "" {
		noHosts, _, e := ParseIPs(noHost)
		if e != nil {
			return nil, nil, e
		}
		if len(noHosts) > 0 {
			var newData []string
			temp := map[string]struct{}{}
			for _, host := range hosts {
				temp[host] = struct{}{}
			}
			for _, host := range noHosts {
				delete(temp, host)
			}
			for host := range temp {
				newData = append(newData, host)
			}
			hosts = newData
			sort.Strings(hosts)
		}
	}
	hosts = arr.SliceRemoveDuplicates(hosts)
	return hosts, domainList, err
}

func ParseIPs(ip string) (hosts, domainList []string, err error) {
	var domain string
	var ips []string
	if strings.Contains(ip, ",") || strings.Contains(ip, "\n") {
		//解析,和\n
		var tmpList []string
		if strings.Contains(ip, ",") {
			tmpList = append(tmpList, strings.Split(ip, ",")...)
		}
		if strings.Contains(ip, "\n") {
			tmpList = append(tmpList, strings.Split(ip, "\n")...)
		}
		//移除重复项
		tmpList = arr.SliceRemoveDuplicates(tmpList)

		//解析多个IP
		for _, item := range tmpList {
			if item == " " || item == "\n" || item == "\r" {
				continue
			}
			item = strings.ReplaceAll(item, " ", "")
			item = strings.ReplaceAll(item, "\n", "")
			item = strings.ReplaceAll(item, "\r", "")
			ips, domain, err = ParseIPone(item)
			if err != nil {
				continue
			}
			if len(ips) != 0 {
				hosts = append(hosts, ips...)
			}
			if domain != "" {
				domainList = append(domainList, domain)
			}
		}

	} else {
		ip = strings.ReplaceAll(ip, " ", "")
		//解析单个IP
		hosts, domain, err = ParseIPone(ip)
		if err != nil {
			return nil, nil, err
		}
		if domain != "" {
			domainList = append(domainList, domain)
		}
		if len(ips) != 0 {
			hosts = append(hosts, ips...)
		}
	}

	//绕过出现报错并且host、domain全部为空进入报错分支
	if err != nil && len(domainList) == 0 && len(hosts) == 0 {
		return nil, nil, err
	}

	return hosts, domainList, nil
}

func ParseIPone(ip string) ([]string, string, error) {
	if len(ip) > 3 {
		switch {
		case strings.Contains(ip, "/24") && !strings.Contains(ip, ":"):
			ipList, err := ParseIPA(ip)
			return ipList, "", err
		case strings.Contains(ip, "/16") && !strings.Contains(ip, ":"):
			ipd, err := ParseIPD(ip)
			return ipd, "", err
		case strings.Contains(ip, "/8") && !strings.Contains(ip, ":"):
			ipe, err := ParseIPE(ip)
			return ipe, "", err
		case strings.Count(ip, "-") == 1 && !IsDomain(ip) && !strings.Contains(ip, ":"):
			ipc, err := ParseIPC(ip)
			return ipc, "", err
		case IsDomain(ip) && !strings.Contains(ip, ":"):
			return []string{}, ip, nil
		default:
			if s := net.ParseIP(ip); s != nil {
				return []string{s.String()}, "", nil
			}
		}
	}

	return []string{}, "", errTips
}

// ParseIPA CIDR IP
func ParseIPA(ip string) ([]string, error) {
	realIP := ip[:len(ip)-3]
	testIP := net.ParseIP(realIP)
	if testIP == nil {
		return nil, errTips
	}
	IPRange := strings.Join(strings.Split(realIP, ".")[0:3], ".")
	var AllIP []string
	for i := 0; i <= 255; i++ {
		AllIP = append(AllIP, IPRange+"."+strconv.Itoa(i))
	}
	return AllIP, nil
}

// ParseIPC Resolving a range of IP,for example: 192.168.111.1-255,192.168.111.1-192.168.112.255
func ParseIPC(ip string) ([]string, error) {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil, errTips
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil, errTips
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil, errTips
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil, errTips
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}
	return AllIP, nil
}

func ParseIPD(ip string) ([]string, error) {
	realIP := ip[:len(ip)-3]
	testIP := net.ParseIP(realIP)
	if testIP == nil {
		return nil, errTips
	}
	IPRange := strings.Join(strings.Split(realIP, ".")[0:2], ".")
	var AllIP []string
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			AllIP = append(AllIP, IPRange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b))
		}
	}
	return AllIP, nil
}

func ParseIPE(ip string) ([]string, error) {
	realIP := ip[:len(ip)-2]
	testIP := net.ParseIP(realIP)
	if testIP == nil {
		return nil, errTips
	}
	IPRange := strings.Join(strings.Split(realIP, ".")[0:1], ".")
	var AllIP []string
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			AllIP = append(AllIP, IPRange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(1))
			AllIP = append(AllIP, IPRange+"."+strconv.Itoa(a)+"."+strconv.Itoa(b)+"."+strconv.Itoa(254))
		}
	}
	return AllIP, nil
}

// IsDomain 判断是不是域名
func IsDomain(s string) bool {
	return regexp.MustCompile("([a-zA-Z-0-9]+|[\u4e00-\u9fa5])+\\.+([a-zA-Z]+|[\u4e00-\u9fa5])").MatchString(s)
}

func IsIPv4(s string) bool {
	pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}

func IsIPv6(s string) bool {
	pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`
	match, _ := regexp.MatchString(pattern, s)
	return match
}
