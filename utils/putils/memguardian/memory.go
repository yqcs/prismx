package memguardian

type SysInfo struct {
	Uptime      int64
	totalRam    uint64
	freeRam     uint64
	SharedRam   uint64
	BufferRam   uint64
	TotalSwap   uint64
	FreeSwap    uint64
	Unit        uint64
	usedPercent float64
}

func (si *SysInfo) TotalRam() uint64 {
	return uint64(si.totalRam) * uint64(si.Unit)
}

func (si *SysInfo) FreeRam() uint64 {
	return uint64(si.freeRam) * uint64(si.Unit)
}

func (si *SysInfo) UsedRam() uint64 {
	return si.TotalRam() - si.FreeRam()
}

func (si *SysInfo) UsedPercent() float64 {
	if si.usedPercent > 0 {
		return si.usedPercent
	}

	return 100 * float64((si.TotalRam()-si.FreeRam())*si.Unit) / float64(si.TotalRam())
}

func GetSysInfo() (*SysInfo, error) {
	return getSysInfo()
}
