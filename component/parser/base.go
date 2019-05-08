package parser

import "fmt"

// Parser 字段分析接口
type Parser interface {
	// Match 配置规则，返回匹配到的值
	Match(rule string) string
	// Match 配置规则，返回匹配到的值，复数
	MatchAll(rule string) []string
}

// Errorf 抛出错误
func Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}
