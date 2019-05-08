package usage

import (
	"github.com/shirou/gopsutil/mem"
)

// MEMStat 内存资源使用情况
type MEMStat struct {
	Total       int
	Used        int
	Free        int
	UsedPercent float64
}

// MEM 获取内存资源使用情况
func MEM() MEMStat {
	s, _ := mem.VirtualMemory()
	return MEMStat{
		Total:       int(s.Total / 1024),
		Used:        int(s.Used / 1024),
		Free:        int(s.Free / 1024),
		UsedPercent: s.UsedPercent,
	}
}
