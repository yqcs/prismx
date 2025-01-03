//go:build !linux

package memguardian

import "github.com/shirou/gopsutil/mem"

// TODO: replace with native syscall
func getSysInfo() (*SysInfo, error) {
	vms, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	si := &SysInfo{
		totalRam:    vms.Total,
		freeRam:     vms.Free,
		SharedRam:   vms.Shared,
		TotalSwap:   vms.SwapTotal,
		FreeSwap:    vms.SwapFree,
		usedPercent: vms.UsedPercent,
	}

	return si, nil
}
