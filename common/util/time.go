package util

import (
	"strconv"
	"time"
)

// SecondToTime 将秒格式化为 几小时:几分钟:几秒
func SecondToTime(second int) string {
	var t []byte
	var hour, minute int

	if second >= 3600 {
		hour = second / 3600
		second = second % 3600
	}

	if second >= 60 {
		minute = second / 60
		second = second % 60
	}

	if hour > 0 {
		if hour < 10 {
			t = append(t, '0')
		}
		t = append(t, strconv.Itoa(hour)...)
		t = append(t, ':')
	}

	if minute < 10 {
		t = append(t, '0')
	}
	t = append(t, strconv.Itoa(minute)...)
	t = append(t, ':')

	if second < 10 {
		t = append(t, '0')
	}
	t = append(t, strconv.Itoa(second)...)

	return string(t)
}

// FriendlyTime 将时间格式化为：
// 几秒前，几分钟前，几小时前，日期 2006-01-02
func FriendlyTime(t time.Time) string {
	now := time.Now()
	if t.Add(time.Second * 60).After(now) {
		return "刚刚"
	}
	if t.Add(time.Second * 3600).After(now) {
		return IntToString(int(now.Unix()-t.Unix())/60) + "分钟前"
	}
	if t.Add(time.Second * 86400).After(now) {
		return IntToString(int(now.Unix()-t.Unix())/3600) + "小时前"
	}
	return t.Format("2006-01-02")
}
