package task

import (
	"fmt"

	"github.com/safeie/spider/component/url"
)

// NewURI 创建一个新URI结构
func (t *Task) NewURI(uri string) *url.URI {
	return url.NewURI(uri)
}

// NewField 创建一个新字段
func (t *Task) NewField(name, alias string) *url.Field {
	return url.NewField(name, alias)
}

// Errorf 抛出错误
func Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}
