package parser

import (
	"regexp"
	"strings"

	"github.com/safeie/spider/common/log"
)

// Regexp 正则解析器
type Regexp struct {
	data string
}

// NewRegexp 创建一个正则匹配解析器
func NewRegexp(body []byte) (*Regexp, error) {
	if body == nil {
		return nil, Errorf("body is empty")
	}
	t := new(Regexp)
	t.data = string(body)
	return t, nil
}

// Match 配置规则，返回匹配到的值
func (t *Regexp) Match(rule string) string {
	if rule == "" {
		return ""
	}
	rule = strings.Replace(regexp.QuoteMeta(rule), "\\(\\*\\)", "(.*)", 1)
	rule = strings.Replace(rule, "\\*", ".*", -1)
	rule = "(?Uis)" + rule
	re, err := regexp.Compile(rule)
	if err != nil {
		log.Errorf("HTML.RegexpMatch Compile Error: %v", err)
		return ""
	}
	matchs := re.FindStringSubmatch(t.data)
	if len(matchs) > 1 {
		return matchs[1]
	}
	return ""
}

// MatchAll 配置规则，返回匹配到的值，复数
// RegexpMatchAll use regexp to match contents, but not real regexp, just support this: <li>(*)</li>
func (t *Regexp) MatchAll(rule string) []string {
	if rule == "" {
		return nil
	}
	rule = strings.Replace(regexp.QuoteMeta(rule), "\\(\\*\\)", "(.*)", 1)
	rule = strings.Replace(rule, "\\*", ".*", -1)
	rule = "(?Uis)" + rule
	re, err := regexp.Compile(rule)
	if err != nil {
		log.Errorf("HTML.RegexpMatch Compile Error: %v", err)
		return nil
	}
	matchs := re.FindAllStringSubmatch(t.data, -1)
	items := make([]string, len(matchs))
	for i, match := range matchs {
		items[i] = match[1]
	}
	return items
}

// Regexp use regexp to match contents, it is a real regexp rule
func (t *Regexp) Regexp(pattern string) []string {
	if pattern == "" {
		return nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re.FindStringSubmatch(t.data)
}
