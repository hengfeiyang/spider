package parser

import (
	"encoding/json"

	"github.com/safeie/spider/common/simplejson"
)

// JSONPath JSON Path解析器
type JSONPath struct {
	parser *simplejson.JSON // JSON解析器
}

// NewJSONPath 创建一个JSON Path解析器
func NewJSONPath(body []byte) (*JSONPath, error) {
	if body == nil {
		return nil, Errorf("body is empty")
	}
	var err error
	t := new(JSONPath)
	t.parser, err = simplejson.New(body)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Match 配置规则，返回匹配到的值
func (t *JSONPath) Match(rule string) interface{} {
	if re := t.parser.Path(rule); re != nil {
		return re.Interface()
	}
	return nil
}

// MatchAll 配置规则，返回匹配到的值，复数
func (t *JSONPath) MatchAll(rule string) []interface{} {
	if re := t.parser.Path(rule); re != nil {
		return re.MustSlice()
	}
	return nil
}

// Marshal 编码JSON
func (t *JSONPath) Marshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
