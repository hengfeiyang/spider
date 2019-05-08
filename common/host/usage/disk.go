package usage

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
)

// DiskStat 磁盘资源使用情况
type DiskStat map[string]DiskItemStat

// DiskItemStat 单个磁盘资源使用情况
type DiskItemStat struct {
	Read  int
	Write int
}

// DISK 获取磁盘资源使用情况，返回由多个磁盘组成的切片
// 计算给定时间(interval)内的硬盘IO平均使用情况，单位：秒，默认值 1，建议值 1
func DISK(interval int) DiskStat {
	if interval == 0 {
		interval = 1
	}
	s1 := CurrentDISKUsage()
	time.Sleep(time.Second * time.Duration(interval))
	s2 := CurrentDISKUsage()
	s := make(DiskStat)
	if len(s1) == 0 || len(s2) == 0 {
		return s
	}
	for dev := range s1 {
		s[dev] = DiskItemStat{
			Read:  (s2[dev].Read - s1[dev].Read),
			Write: (s2[dev].Write - s1[dev].Write),
		}
	}
	return s
}

// CurrentDISKUsage 获取当前磁盘统计
func CurrentDISKUsage() DiskStat {
	stat := make(DiskStat)
	s, err := disk.IOCounters()
	if err != nil {
		return stat
	}
	for _, dev := range s {
		// 过滤重复的设备
		if strings.HasPrefix(dev.Name, "dm-") {
			continue
		}
		// 过滤分区数据
		if len(dev.Name) > 1 {
			m := dev.Name[len(dev.Name)-1]
			// 48 -> 0, 57 -> 9
			if m >= 48 && m <= 57 {
				continue
			}
		}
		stat[dev.Name] = DiskItemStat{
			Read:  int(dev.ReadBytes / 1024),
			Write: int(dev.WriteBytes / 1024),
		}
	}
	return stat
}
