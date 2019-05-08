package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func getTestData() map[string]interface{} {
	m := make(map[string]interface{}, 9)
	m["projectid"] = make([]string, 1)
	m["projectid"].([]string)[0] = "123"
	m["key[]"] = make([]string, 2)
	m["key[]"].([]string)[0] = "123"
	m["key[]"].([]string)[1] = "abc"
	m["time"] = make([]string, 1)
	m["time"].([]string)[0] = "123"
	m["sign"] = make([]string, 1)
	m["sign"].([]string)[0] = "123"
	m["file"] = make([]string, 1)
	m["file"].([]string)[0] = "123"
	m["filesign"] = make([]string, 1)
	m["filesign"].([]string)[0] = "abcdeadsfa"
	m["special"] = make([]string, 1)
	m["special"].([]string)[0] = "123 abc#123+abc%123&abc*123~abc"
	return m
}

func TestUrlParamSort1(t *testing.T) {
	Convey("测试URL参数类型排序", t, func() {
		m := getTestData()
		nm := urlParamSort(m)
		So(nm[1].Key, ShouldEqual, "filesign")
		So(nm[6].Key, ShouldEqual, "time")
	})
}

func TestHTTPBuildQuery1(t *testing.T) {
	Convey("测试HTTPBuildQuery", t, func() {
		m := getTestData()
		s := HTTPBuildQuery(m, QUERY_RFC3986)
		So(s, ShouldEqual, "file=123&filesign=abcdeadsfa&key%5B%5D=123&key%5B%5D=abc&projectid=123&sign=123&special=123%20abc%23123%2Babc%25123%26abc%2A123~abc&time=123")
	})
}
