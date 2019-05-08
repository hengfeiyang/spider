package parser

import (
	"regexp"
	"strings"

	"github.com/safeie/spider/common/log"
)

// Substring 字符串匹配
type Substring struct {
	data string
}

// NewSubstring 创建一个字符串匹配解析器
func NewSubstring(body []byte) (*Substring, error) {
	if body == nil {
		return nil, Errorf("body is empty")
	}
	t := new(Substring)
	t.data = string(body)
	return t, nil
}

// Match 配置规则，返回匹配到的值
func (t *Substring) Match(rule string) string {
	patternArr := strings.Split(rule, "(*)")
	if len(patternArr) != 2 {
		log.Debugln("Substring.Match: rule format error, must contains/split (*)")
	}
	return t.StrMatch(patternArr[0], patternArr[1])
}

// MatchAll 配置规则，返回匹配到的值，复数
func (t *Substring) MatchAll(rule string) []string {
	patternArr := strings.Split(rule, "(*)")
	if len(patternArr) != 2 {
		log.Debugln("Substring.MatchAll: rule format error, must contains/split (*)")
	}
	return t.StrMatchAll(patternArr[0], patternArr[1])
}

// StrMatch use begin and end code to match content
func (t *Substring) StrMatch(start, end string) string {
	if start == "" || end == "" {
		return ""
	}

	pos1 := strings.Index(t.data, start)
	if pos1 == -1 {
		return ""
	}
	pos1 += len(start)

	pos2 := strings.Index(t.data[pos1:], end)
	if pos2 == -1 {
		return ""
	}
	pos2 += pos1
	return t.data[pos1:pos2]
}

// StrMatchAll use begin and end code to match all contents
func (t *Substring) StrMatchAll(start, end string) []string {
	return t.RegexpMatchAll(start + "(*)" + end)
}

// RegexpMatchAll use regexp to match contents, but not real regexp, just support this: <li>(*)</li>
func (t *Substring) RegexpMatchAll(pattern string) []string {
	if pattern == "" {
		return nil
	}
	pattern = strings.Replace(regexp.QuoteMeta(pattern), "\\(\\*\\)", "(.*)", 1)
	pattern = strings.Replace(pattern, "\\*", ".*", -1)
	pattern = "(?Uis)" + pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Debugln("Substring.RegexpMatch Compile Error:", err)
		return nil
	}
	matchs := re.FindAllStringSubmatch(t.data, -1)
	items := make([]string, len(matchs))
	for i, match := range matchs {
		items[i] = match[1]
	}
	return items
}
