package parse

import (
	"prismx_cli/utils/arr"
	"sort"
	"strconv"
	"strings"
)

func ParsePort(ports string) (scanPorts []int) {
	if ports == "" {
		return
	}
	slices := strings.Split(ports, ",")
	for _, port := range slices {
		port = strings.Trim(port, " ")
		upper := port
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}

			startPort, _ := strconv.Atoi(ranges[0])
			endPort, _ := strconv.Atoi(ranges[1])
			if startPort < endPort {
				port = ranges[0]
				upper = ranges[1]
			} else {
				port = ranges[1]
				upper = ranges[0]
			}

		}
		start, _ := strconv.Atoi(port)
		end, _ := strconv.Atoi(upper)
		for i := start; i <= end; i++ {
			scanPorts = append(scanPorts, i)
		}
	}
	scanPorts = arr.IntSliceRemoveDuplicates(scanPorts)
	return scanPorts
}

// GetScanPort 获取待扫描的端口
func GetScanPort(p, bp string) []int {
	probePorts := ParsePort(p)
	noPorts := ParsePort(bp)
	if len(noPorts) > 0 {
		var newData []int
		temp := map[int]struct{}{}
		for _, port := range probePorts {
			temp[port] = struct{}{}
		}
		for _, port := range noPorts {
			delete(temp, port)
		}
		for port := range temp {
			newData = append(newData, port)
		}
		probePorts = newData
		sort.Ints(probePorts)
	}
	return probePorts
}
