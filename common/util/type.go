package util

import (
	"strconv"
)

// Int2String format int to string
func Int2String(i int) string {
	s := strconv.Itoa(i)
	return s
}

// String2Int format string to int
func String2Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
