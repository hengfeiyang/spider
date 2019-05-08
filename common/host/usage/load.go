package usage

import (
	"github.com/shirou/gopsutil/load"
)

// LoadStat 系统负载情况
type LoadStat struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// Load 获取系统负载情况
func Load() LoadStat {
	s, _ := load.Avg()
	return LoadStat{
		Load1:  s.Load1,
		Load5:  s.Load5,
		Load15: s.Load15,
	}
}
