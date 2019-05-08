package usage

import (
	"time"

	"github.com/safeie/spider/common/util"
	"github.com/shirou/gopsutil/cpu"
)

// CPUStat CPU资源使用情况
type CPUStat []int

// CPU 获取CPU资源使用情况，返回由多核心CPU资源组成的切片
// 计算给定时间(interval)内的CPU平均使用情况，单位：秒，默认值 1，建议值 1
func CPU(interval int) CPUStat {
	if interval == 0 {
		interval = 1
	}
	stat := make(CPUStat, 0)
	s, err := cpu.Percent(time.Second*time.Duration(interval), true)
	if err != nil {
		return stat
	}
	for _, c := range s {
		stat = append(stat, util.MathRound(c))
	}
	return stat
}
