package usage

import (
	"time"

	"github.com/shirou/gopsutil/net"
)

// NetStat 网络资源使用情况
type NetStat map[string]NetItemStat

// NetItemStat 单个网卡的资源使用情况
type NetItemStat struct {
	In  int
	Out int
}

// NET 获取网络资源使用情况，返回由多个网卡组成的切片
// 计算给定时间(interval)内的网络IO平均使用情况，单位：秒，默认值 1，建议值 1
func NET(interval int) NetStat {
	if interval == 0 {
		interval = 1
	}
	s1 := CurrentNetUsage()
	time.Sleep(time.Second * time.Duration(interval))
	s2 := CurrentNetUsage()
	s := make(NetStat)
	if len(s1) == 0 || len(s2) == 0 {
		return s
	}
	for itr := range s1 {
		s[itr] = NetItemStat{
			In:  (s2[itr].In - s1[itr].In),
			Out: (s2[itr].Out - s1[itr].Out),
		}
	}
	return s
}

// CurrentNetUsage 获取当前网卡统计
func CurrentNetUsage() NetStat {
	stat := make(NetStat)
	s, err := net.IOCounters(true)
	if err != nil {
		return stat
	}
	for _, itr := range s {
		stat[itr.Name] = NetItemStat{
			In:  int(itr.BytesRecv / 1024),
			Out: int(itr.BytesSent / 1024),
		}
	}
	return stat
}
