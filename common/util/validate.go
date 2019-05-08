package util

import (
	"regexp"
)

func Validate(str, category string) bool {
	var matched bool
	if str == "" {
		return matched
	}
	switch category {
	case "mobile":
		matched, _ = regexp.MatchString(`^1[3|4|5|7|8]\d{9}$`, str)
	case "number":
		matched, _ = regexp.MatchString(`\d+`, str)
	case "email":
		matched, _ = regexp.MatchString(`\w+([-+.]\w+)*@\w+([-+.]\w+)*.\w+([-+.]\w+)*`, str)
	case "url":
		matched, _ = regexp.MatchString(`http(s)?\://[a-z0-9-.]+/?.*`, str)
	case "ipv4":
		matched, _ = regexp.MatchString(`((1-9)|[12][0-9]|[12][0-5][0-5]).((0-9)|[12][0-9]|[12][0-5][0-5]).((0-9)|[12][0-9]|[12][0-5][0-5]).((1-9)|[12][0-9]|[12][0-5][0-5])`, str)
	}
	return matched
}
