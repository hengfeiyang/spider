package util

import (
	"net/url"
	"sort"
	"strings"
)

const (
	// QUERY_RFC1738 则编码将会以 » RFC 1738 标准和 application/x-www-form-urlencoded 媒体类型进行编码，空格会被编码成加号（+）。
	QUERY_RFC1738 int = iota
	// QUERY_RFC3986 将根据 » RFC 3986 编码，空格会被百分号编码（%20）。
	QUERY_RFC3986
)

type urlValue []urlValueItem

type urlValueItem struct {
	Key string
	Val interface{}
}

func (t urlValue) Len() int {
	return len(t)
}

func (t urlValue) Less(j, k int) bool {
	return t[j].Key < t[k].Key
}

func (t urlValue) Swap(j, k int) {
	t[j], t[k] = t[k], t[j]
}

func (t urlValue) Sort() {
	sort.Sort(t)
}

func urlParamSort(mapItems map[string]interface{}) []urlValueItem {
	ns := make(urlValue, 0, len(mapItems))
	for k, v := range mapItems {
		ns = append(ns, urlValueItem{Key: k, Val: v})
	}
	ns.Sort()

	return ns
}

// HTTPBuildQuery 模仿php的http_build_query构建字符串
// JAVA String resultMD5 = MD5.md5(result.toString().replace("*", "%2A").replace("%7E", "~").replace("+", "%20"));
func HTTPBuildQuery(values map[string]interface{}, encType int) string {
	if encType == 0 {
		encType = QUERY_RFC3986
	}
	s := make([]string, 0, len(values))
	nv := urlParamSort(values)
	for _, item := range nv {
		key := url.QueryEscape(item.Key)
		if val, ok := item.Val.([]string); ok {
			for _, v := range val {
				s = append(s, key+"="+url.QueryEscape(v))
			}
		} else {
			s = append(s, key+"="+url.QueryEscape(item.Val.(string)))
		}
	}
	ss := strings.Join(s, "&")
	if encType == QUERY_RFC3986 {
		ss = strings.Replace(ss, "+", "%20", -1)
	}
	return ss
}
