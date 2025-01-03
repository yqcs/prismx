//go:build linux

package memguardian

import "syscall"

func getSysInfo() (*SysInfo, error) {
	var sysInfo syscall.Sysinfo_t
	err := syscall.Sysinfo(&sysInfo)
	if err != nil {
		return nil, err
	}

	si := &SysInfo{
		Uptime:    int64(sysInfo.Uptime),
		totalRam:  uint64(sysInfo.Totalram),
		freeRam:   uint64(sysInfo.Freeram),
		SharedRam: uint64(sysInfo.Freeram),
		BufferRam: uint64(sysInfo.Bufferram),
		TotalSwap: uint64(sysInfo.Totalswap),
		FreeSwap:  uint64(sysInfo.Freeswap),
		Unit:      uint64(sysInfo.Unit),
	}

	return si, nil
}
