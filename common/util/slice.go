package util

import (
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// SliceIntEqual 判断两个数字切片中的内容是否相同，忽略元素的排序
// [1,2,3] == [3,2,1]
func SliceIntEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SliceIntToString 将数字切片转换为字符串切片
// [11,22,33] 转换为 ["11","22","33"]
func SliceIntToString(a []int) []string {
	str := make([]string, len(a))
	for k := range a {
		str[k] = strconv.Itoa(a[k])
	}
	return str
}

// SliceIntRand 将切片顺序打乱随机返回
func SliceIntRand(a []int) []int {
	if len(a) <= 1 {
		return a
	}
	aNew := make([]int, len(a))
	rand.Seed(time.Now().UnixNano())
	n := len(a)
	for i := 0; i < n; i++ {
		m := len(a)
		if m == 1 {
			aNew[i] = a[0]
			break
		}
		k := rand.Intn(m)
		aNew[i] = a[k]
		b := a[0:k]
		b = append(b, a[k+1:]...)
		a = b
	}
	return aNew
}

// InSlice 判断元素s是否在slice si中出现过,返回 bool
func InSlice(si interface{}, s interface{}, t string) bool {
	if reflect.ValueOf(si).Len() == 0 {
		return false
	}
	switch t {
	case "string":
		for _, v := range si.([]string) {
			if v == s.(string) {
				return true
			}
		}
	case "Time", "time":
		for _, v := range si.([]time.Time) {
			if v == s.(time.Time) {
				return true
			}
		}
	case "int":
		for _, v := range si.([]int) {
			if v == s.(int) {
				return true
			}
		}
	case "int8":
		for _, v := range si.([]int8) {
			if v == s.(int8) {
				return true
			}
		}
	case "int16":
		for _, v := range si.([]int16) {
			if v == s.(int16) {
				return true
			}
		}
	case "int32":
		for _, v := range si.([]int32) {
			if v == s.(int32) {
				return true
			}
		}
	case "int64":
		for _, v := range si.([]int64) {
			if v == s.(int64) {
				return true
			}
		}
	case "uint":
		for _, v := range si.([]uint) {
			if v == s.(uint) {
				return true
			}
		}
	case "uint8":
		for _, v := range si.([]uint8) {
			if v == s.(uint8) {
				return true
			}
		}
	case "uint16":
		for _, v := range si.([]uint16) {
			if v == s.(uint16) {
				return true
			}
		}
	case "uint32":
		for _, v := range si.([]uint32) {
			if v == s.(uint32) {
				return true
			}
		}
	case "uint64":
		for _, v := range si.([]uint64) {
			if v == s.(uint64) {
				return true
			}
		}
	case "float32":
		for _, v := range si.([]float32) {
			if v == s.(float32) {
				return true
			}
		}
	case "float64":
		for _, v := range si.([]float64) {
			if v == s.(float64) {
				return true
			}
		}
	default:
		return false
	}
	return false
}

// SliceIntDiff 获取在数字切片 a 中但不在数字切片 b 中的差集
func SliceIntDiff(a []int, b []int) []int {
	var diff []int

	for _, s1 := range a {
		found := false
		for _, s2 := range b {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// SliceStringDiff 获取在字符串切片 a 中但不在字符串切片 b 中的差集
func SliceStringDiff(a []string, b []string) []string {
	var diff []string

	for _, s1 := range a {
		found := false
		for _, s2 := range b {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

// MixedToSliceInt 将混合类型转为 数字slice
func MixedToSliceInt(v interface{}) []int {
	var result []int
	var _id int
	if _, ok := v.(string); ok {
		_id, _ = strconv.Atoi(v.(string))
		if _id > 0 {
			result = append(result, _id)
		}
	} else if _, ok := v.([]string); ok {
		for i := range v.([]string) {
			_id, _ = strconv.Atoi(v.([]string)[i])
			if _id > 0 {
				result = append(result, _id)
			}
		}
	} else if _, ok := v.(int); ok {
		result = append(result, v.(int))
	} else if _, ok := v.([]int); ok {
		result = v.([]int)
	}
	return result
}

// MixedToSliceString 将混合类型转为 字符串slice
func MixedToSliceString(v interface{}) []string {
	var result []string
	if _, ok := v.(string); ok {
		result = append(result, v.(string))
	} else if _, ok := v.([]string); ok {
		result = v.([]string)
	} else if _, ok := v.(int); ok {
		result = append(result, strconv.Itoa(v.(int)))
	} else if _, ok := v.([]int); ok {
		for i := range v.([]int) {
			result = append(result, strconv.Itoa(v.([]int)[i]))
		}
	}
	return result
}
