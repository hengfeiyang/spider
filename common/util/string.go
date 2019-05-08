package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	STR_PAD_LEFT int = iota
	STR_PAD_RIGHT
	STR_PAD_BOTH
)

// ToString 将任意一个类型转换为字符串
func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

// StringToInt 将字符串转为int
func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// StringToInt64 将字符串转为int64
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// IntToString 将数字转换为字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// CamelCase 将一个字符串转为大驼峰命名
func CamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	ss := strings.Split(s, " ")
	for k, v := range ss {
		ss[k] = strings.Title(v)
	}
	return strings.Join(ss, "")
}

// Link: https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// CamelCaseInitialism 将一个字符串转为大驼峰命名，强制首字母缩写命名规范
func CamelCaseInitialism(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	ss := strings.Split(s, " ")
	var sm bool
	var uv string
	for k, v := range ss {
		sm = false
		uv = strings.ToUpper(v)
		for m := range commonInitialisms {
			if uv == m {
				ss[k] = m
				sm = true
				break
			}
		}
		if !sm {
			ss[k] = strings.Title(v)
		}
	}
	return strings.Join(ss, "")
}

// IsNumeric 判断是否是纯数字，空字符串不是数字，0才是
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		fmt.Println(v)
		if v < 48 || v > 57 {
			return false
		}
	}
	return true
}

// StrPad 使用另一个字符串填充字符串为指定长度
func StrPad(v interface{}, length int, pad string, padType int) string {
	s := fmt.Sprintf("%v", v)
	if len(s) >= length {
		return s
	}
	var pos int
	switch padType {
	case STR_PAD_BOTH:
		left := true
		for len(s) < length {
			// first left, then right
			if left {
				s = pad + s
				if len(s) > length {
					pos = len(s) - length
					return s[pos:]
				}
				left = false
			} else {
				s = s + pad
				if len(s) > length {
					return s[:length]
				}
				left = true
			}
		}
	case STR_PAD_RIGHT:
		for len(s) < length {
			s = s + pad
			if len(s) > length {
				return s[:length]
			}
		}
	default:
		for len(s) < length {
			s = pad + s
			if len(s) > length {
				pos = len(s) - length
				return s[pos:]
			}
		}

	}

	return s
}

// StrNatCut 字符串截取，中文算一个 英文算两个
func StrNatCut(s string, length int, dot ...string) string {
	if len(s) < length {
		return s
	}
	if dot == nil {
		dot = []string{"..."}
	}

	n := 0
	l := len(s)
	for i, b := range s {
		if n >= length*2 {
			if l > i {
				return s[:i] + strings.Join(dot, "")
			}
			return s[:i]
		}
		if b < 128 {
			n++
		} else {
			n += 2
		}
	}

	return s
}

// Concat 连接字符串
func Concat(args ...string) string {
	var buf bytes.Buffer
	for i := range args {
		buf.WriteString(args[i])
	}
	return buf.String()
}
