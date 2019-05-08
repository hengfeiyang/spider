package simplejson

import (
	"encoding/json"
	"strconv"
)

// JSON json解析器
type JSON struct {
	data interface{}
}

// New 创建新的JSON解析器
func New(data []byte) (*JSON, error) {
	j := new(JSON)
	if data != nil {
		err := j.Unmarshal(data)
		if err != nil {
			return nil, err
		}
	}
	return j, nil
}

// Unmarshal 解析JSON数据到解析器
func (t *JSON) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &t.data)
}

// Marshal 将JSON解析器中的数据，编码成JSON字符串
func (t *JSON) Marshal() ([]byte, error) {
	return json.Marshal(t.data)
}

// Path 解析Path路径
// $.data
// $.data.name
// $.data[0].name
// $[0].name
func (t *JSON) Path(path string) *JSON {
	if len(path) < 3 {
		return nil
	}
	if path[0] != '$' {
		return nil
	}
	path = path[1:]
	var j = t
	var key []rune
	var ikey int
	var lastChar rune
	for _, k := range path {
		if k == '.' || k == '[' {
			if len(key) == 0 || lastChar == ']' {
				continue
			}
			j = t.Get(string(key))
			key = make([]rune, 0)
			lastChar = k
			continue
		}
		if k == ']' {
			if len(key) == 0 {
				return nil
			}
			ikey, _ = strconv.Atoi(string(key))
			j = j.GetIndex(ikey)
			key = make([]rune, 0)
			lastChar = k
			continue
		}
		key = append(key, k)
		lastChar = k
	}
	if len(key) > 0 {
		j = j.Get(string(key))
	}
	return j
}

// Get 获取一个结构字段
func (t *JSON) Get(path string) *JSON {
	if j, ok := t.data.(map[string]interface{}); ok {
		s, err := json.Marshal(j[path])
		if err == nil {
			j, _ := New(s)
			return j
		}
	}
	return nil
}

// GetIndex 获取数组中的一个值
func (t *JSON) GetIndex(i int) *JSON {
	if j, ok := t.data.([]interface{}); ok {
		if s, err := json.Marshal(j[i]); err == nil {
			if j, err := New(s); err == nil {
				return j
			}
		}
	}
	return nil
}

// Interface 输出值为 原始类型
func (t *JSON) Interface() interface{} {
	return t.data
}

// String 实现 Stringer 接口
func (t *JSON) String() string {
	if j, ok := t.data.(string); ok {
		return j
	}
	if j, err := t.Marshal(); err == nil {
		return string(j)
	}
	return ""
}

// MustString 输出值为 字符串
func (t *JSON) MustString() string {
	if j, ok := t.data.(string); ok {
		return j
	}
	return ""
}

// MustInt 输出值为 数字
func (t *JSON) MustInt() int {
	j := t.MustFloat64()
	return int(j)
}

// MustFloat64 输出值为 浮点数
func (t *JSON) MustFloat64() float64 {
	if j, ok := t.data.(float64); ok {
		return j
	}
	return 0
}

// MustSlice 输出值为 数组/切片
func (t *JSON) MustSlice() []interface{} {
	if j, ok := t.data.([]interface{}); ok {
		return j
	}
	return nil
}

// MustMap 输出值为 字典/Map
func (t *JSON) MustMap() map[string]interface{} {
	if j, ok := t.data.(map[string]interface{}); ok {
		return j
	}
	return nil
}
