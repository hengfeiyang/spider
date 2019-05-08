package util

import "strings"

// VersionDiff 如果v1大于等于v2返回true, 否则返回 false
func VersionDiff(v1, v2 string) bool {
	v1a := strings.Split(v1, ".")
	v2a := strings.Split(v2, ".")
	if d := len(v1a) - len(v2a); d > 0 {
		for i := 0; i < d; i++ {
			v2a = append(v2a, "0")
		}
	}
	for i := range v1a {
		if v1a[i] < v2a[i] {
			return false
		}
	}
	return true
}
