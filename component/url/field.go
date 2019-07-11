package url

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/safeie/spider/common/log"
	"github.com/safeie/spider/common/util"
	"github.com/safeie/spider/component/parser"
)

const (
	MatchTypeSelector  = iota // html Selector匹配
	MatchTypeSubString        // 字符串截取
	MatchTypeRegexp           // 正则表达式
	MatchTypeJSONPath         // json path匹配
)

const (
	SourceTypeContext = iota // 当前内容
	SourceTypeAttach         // URL附件字段
)

// FieldFilterFunc 字段过滤方法
type FieldFilterFunc func(f *Field)

// Field 字段
type Field struct {
	Name        string              // 字段名
	Alias       string              // 别名（中文）
	Remote      *Remote             // 远程获取字段附属页面
	value       interface{}         // 字段值
	sourceType  int                 // 字段来源，默认 当前页面中，可选，附加字段，远程字段 page,attach,remote
	matchType   int                 // 匹配类型
	matchRule   string              // 匹配规则
	expand      bool                // 展开单个复数字段，即：只有一个孩子字段且该字段为数组时，展开该字段为多条数据
	fixURL      bool                // 是否修复URL，修复可能的相对路径
	fixed       bool                // 是否已经修复过URL
	repeat      bool                // 是否重复
	repeatValue []map[string]*Field // 重复字段值
	children    []*Field            // 孩子字段
	filters     []FieldFilterFunc   // 过滤函数集，优先于全局规则过滤器执行

}

// NewField 创建一个新字段
func NewField(name, alias string) *Field {
	f := new(Field)
	f.Name = name
	f.Alias = alias
	f.children = make([]*Field, 0)
	return f
}

// Copy 复制出一个新字段
func (f *Field) Copy() *Field {
	n := new(Field)
	n.Name = f.Name
	n.Alias = f.Alias
	n.value = f.value
	if f.Remote != nil {
		n.Remote = &Remote{
			url:         f.Remote.url,
			pageType:    f.Remote.pageType,
			fetchOption: f.Remote.fetchOption,
			logger:      f.Remote.logger,
		}
	}
	n.sourceType = f.sourceType
	n.matchType = f.matchType
	n.matchRule = f.matchRule
	n.expand = f.expand
	n.repeat = f.repeat
	n.fixURL = f.fixURL
	for i := range f.children {
		n.children = append(n.children, f.children[i].Copy())
	}
	n.filters = f.filters
	return n
}

// Children 获取带值的子字段
func (f *Field) Children() []map[string]*Field {
	return f.repeatValue
}

// Value 获取字段的值，递归处理
func (f *Field) Value() interface{} {
	if f.children == nil {
		return f.value
	}
	if f.repeat {
		value := make([]map[string]interface{}, 0)
		for _, repeat := range f.repeatValue {
			cvalue := make(map[string]interface{})
			for key, cf := range repeat {
				cvalue[key] = cf.Value()
			}
			value = append(value, cvalue)
		}
		return value
	}
	value := make(map[string]interface{})
	for _, cf := range f.children {
		value[cf.Name] = cf.Value()
	}
	if f.expand && len(value) == 1 {
		for k := range value {
			return value[k]
		}
	}
	return value
}

// String 获取字段值的string类型
func (f *Field) String() string {
	return string(f.Bytes())
}

// Bytes 获取字段值的byte类型
func (f *Field) Bytes() []byte {
	if v, ok := f.value.(string); ok {
		return []byte(v)
	}
	if v, ok := f.value.([]byte); ok {
		return v
	}
	return nil
}

// Replace 替换字段中值出现的字符串，支持 正则表达式
// 支持，伪正则，* 号 = .*
func (f *Field) Replace(old, new string) {
	if old == "" {
		return
	}
	// 非字符串类型，不能替换
	var value string
	if value = f.String(); len(value) == 0 {
		return
	}
	old = strings.Replace(regexp.QuoteMeta(old), "\\(\\*\\)", "(.*)", 1)
	old = strings.Replace(old, "\\*", ".*", -1)
	old = "(?Uis)" + old
	re, err := regexp.Compile(old)
	if err != nil {
		log.Errorln("Field.Replace Compile Error:", err)
		return
	}
	f.SetValue(re.ReplaceAllString(value, new))
}

// Remove 删除字段中值出现的字符串，支持 正则表达式
// 支持，伪正则，* 号 = .*
func (f *Field) Remove(rule string) {
	f.Replace(rule, "")
}

// RemoveString 删除字段中值出现的字符串
func (f *Field) RemoveString(str string) {
	v := f.String()
	if len(v) == 0 {
		return
	}
	f.SetValue(strings.ReplaceAll(v, str, ""))
}

// RemoveHTMLTags 过滤字段值中出现的HTML标签
func (f *Field) RemoveHTMLTags(tags ...string) {
	v := f.Bytes()
	if len(v) == 0 {
		return
	}
	v = util.StripTags(v, tags...)
	// 注意，要以string设置回去，因为 其他过滤器都依赖string类型
	f.SetValue(string(v))
}

// Trim 过滤字段值中两端的字符
func (f *Field) Trim(cutset string) {
	v := f.String()
	if len(v) == 0 {
		return
	}
	f.SetValue(strings.Trim(v, cutset))
}

// TrimSpace 过滤字段值中的空白
func (f *Field) TrimSpace() {
	v := f.String()
	if len(v) == 0 {
		return
	}
	f.SetValue(strings.TrimSpace(v))
}

// SetValue 设置字段值
func (f *Field) SetValue(v interface{}) *Field {
	f.value = v
	return f
}

// SetMatchRule 设置匹配规则
func (f *Field) SetMatchRule(matchType int, matchRule string) *Field {
	f.matchType = matchType
	f.matchRule = matchRule
	return f
}

// SetSourceType 设置来源类型，默认 页面本身
// 支持：页面 page，远程 remote，附加字段 attach
func (f *Field) SetSourceType(v int) *Field {
	f.sourceType = v
	return f
}

// SetExpand 展开单个复数字段，即：只有一个孩子字段且该字段为数组时，展开该字段为多条数据
func (f *Field) SetExpand(v bool) *Field {
	f.expand = v
	return f
}

// SetRemote 设置远程获取属性
func (f *Field) SetRemote(r *Remote) *Field {
	f.Remote = r
	return f
}

// SetRepeat 设置是否是重复项
func (f *Field) SetRepeat(v bool) *Field {
	f.repeat = v
	return f
}

// SetFixURL 设置是否是否修复URL，修复可能的相对路径
func (f *Field) SetFixURL(v bool) *Field {
	f.fixURL = v
	return f
}

// SetFilterFunc 设置过滤函数，仅对当前字段有效
func (f *Field) SetFilterFunc(ff ...FieldFilterFunc) *Field {
	f.filters = append(f.filters, ff...)
	return f
}

// SetChildren 设置子字段，导出字段时，作为该字段的子集
func (f *Field) SetChildren(fc ...*Field) *Field {
	f.children = append(f.children, fc...)
	return f
}

// Fetch 获取字段值
func (f *Field) Fetch(u *URI) error {
	// 附件字段
	if f.sourceType == SourceTypeAttach {
		f.value = u.Get(f.Name)
		return nil
	}

	if f.matchRule == "" {
		return fmt.Errorf("字段提取规则为空")
	}
	var err error
	switch f.matchType {
	case MatchTypeSelector:
		err = f.fetchSelector(u)
	case MatchTypeSubString:
		err = f.fetchSubstring(u)
	case MatchTypeRegexp:
		err = f.fetchRegexp(u)
	case MatchTypeJSONPath:
		err = f.fetchJSON(u)
	default:
		err = fmt.Errorf("不支持的字段提取类型")
	}
	// 执行字段过滤器
	if err == nil {
		for _, fn := range f.filters {
			fn(f)
		}
	}
	// 修复URL
	if f.fixURL {
		err = f.fixValueURL(u)
	}

	return err
}

// fetchJSON 解析JSON
func (f *Field) fetchJSON(u *URI) error {
	var err error
	if u.Parser.JSON == nil {
		if u.Parser.JSON, err = parser.NewJSONPath(u.Body); err != nil {
			return err
		}
	}

	if f.repeat {
		f.value = u.Parser.JSON.MatchAll(f.matchRule)
	} else {
		f.value = u.Parser.JSON.Match(f.matchRule)
	}

	if f.children != nil {
		furi := u.Copy()
		if f.repeat {
			value, ok := f.value.([]interface{})
			if !ok {
				return nil
			}
			for i := range value {
				if f.Remote != nil {
					if furi, err = f.Remote.FetchURI(string(u.Parser.JSON.Marshal(value[i]))); err != nil {
						return err
					}
				} else {
					furi.ResetBody(u.Parser.JSON.Marshal(value[i]))
				}
				repeatValue := make(map[string]*Field)
				for j := range f.children {
					cf := f.children[j].Copy()
					if err = cf.Fetch(furi); err != nil {
						return err
					}
					repeatValue[cf.Name] = cf
				}
				f.repeatValue = append(f.repeatValue, repeatValue)
			}
		} else {
			if f.Remote != nil {
				if furi, err = f.Remote.FetchURI(string(u.Parser.JSON.Marshal(f.value))); err != nil {
					return err
				}
			} else {
				furi.ResetBody(u.Parser.JSON.Marshal(f.value))
			}
			for _, cf := range f.children {
				if err = cf.Fetch(furi); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// fetchSelector 解析html dom selector
func (f *Field) fetchSelector(u *URI) error {
	var err error
	if u.Parser.HTMLDom == nil {
		if u.Parser.HTMLDom, err = parser.NewHTMLDom(u.Body); err != nil {
			return err
		}
	}

	if f.repeat {
		f.value = u.Parser.HTMLDom.MatchAll(f.matchRule)
	} else {
		f.value = u.Parser.HTMLDom.Match(f.matchRule)
	}

	if f.children != nil {
		furi := u.Copy()
		if f.repeat {
			value, ok := f.value.([]string)
			if !ok {
				return nil
			}
			for i := range value {
				if f.Remote != nil {
					if furi, err = f.Remote.FetchURI(value[i]); err != nil {
						return err
					}
				} else {
					furi.ResetBody([]byte(value[i]))
				}
				repeatValue := make(map[string]*Field)
				for j := range f.children {
					cf := f.children[j].Copy()
					if err = cf.Fetch(furi); err != nil {
						return err
					}
					repeatValue[cf.Name] = cf
				}
				f.repeatValue = append(f.repeatValue, repeatValue)
			}
		} else {
			if f.Remote != nil {
				if furi, err = f.Remote.FetchURI(f.value.(string)); err != nil {
					return err
				}
			} else {
				furi.ResetBody([]byte(f.value.(string)))
			}
			for _, cf := range f.children {
				if err = cf.Fetch(furi); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// fetchRegexp 正则解析
func (f *Field) fetchRegexp(u *URI) error {
	var err error
	if u.Parser.Regexp == nil {
		if u.Parser.Regexp, err = parser.NewRegexp(u.Body); err != nil {
			return err
		}
	}

	if f.repeat {
		f.value = u.Parser.Regexp.MatchAll(f.matchRule)
	} else {
		f.value = u.Parser.Regexp.Match(f.matchRule)
	}

	if f.children != nil {
		furi := u.Copy()
		if f.repeat {
			value, ok := f.value.([]string)
			if !ok {
				return nil
			}
			for i := range value {
				if f.Remote != nil {
					if furi, err = f.Remote.FetchURI(value[i]); err != nil {
						return err
					}
				} else {
					furi.ResetBody([]byte(value[i]))
				}
				repeatValue := make(map[string]*Field)
				for j := range f.children {
					cf := f.children[j].Copy()
					if err = cf.Fetch(furi); err != nil {
						return err
					}
					repeatValue[cf.Name] = cf
				}
				f.repeatValue = append(f.repeatValue, repeatValue)
			}
		} else {
			if f.Remote != nil {
				if furi, err = f.Remote.FetchURI(f.value.(string)); err != nil {
					return err
				}
			} else {
				furi.ResetBody([]byte(f.value.(string)))
			}
			for _, cf := range f.children {
				if err = cf.Fetch(furi); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// fetchSubstring 字符串匹配
func (f *Field) fetchSubstring(u *URI) error {
	var err error
	if u.Parser.Substring == nil {
		if u.Parser.Substring, err = parser.NewSubstring(u.Body); err != nil {
			return err
		}
	}

	if f.repeat {
		f.value = u.Parser.Substring.MatchAll(f.matchRule)
	} else {
		f.value = u.Parser.Substring.Match(f.matchRule)
	}

	if f.children != nil {
		furi := u.Copy()
		if f.repeat {
			value, ok := f.value.([]string)
			if !ok {
				return nil
			}
			for i := range value {
				if f.Remote != nil {
					if furi, err = f.Remote.FetchURI(value[i]); err != nil {
						return err
					}
				} else {
					furi.ResetBody([]byte(value[i]))
				}
				repeatValue := make(map[string]*Field)
				for j := range f.children {
					cf := f.children[j].Copy()
					if err = cf.Fetch(furi); err != nil {
						return err
					}
					repeatValue[cf.Name] = cf
				}
				f.repeatValue = append(f.repeatValue, repeatValue)
			}
		} else {
			if f.Remote != nil {
				if furi, err = f.Remote.FetchURI(f.value.(string)); err != nil {
					return err
				}
			} else {
				furi.ResetBody([]byte(f.value.(string)))
			}
			for _, cf := range f.children {
				if err = cf.Fetch(furi); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// fixValueURL 修复内容中的链接，图片路径为完整姿势
func (f *Field) fixValueURL(u *URI) error {
	if f.fixed {
		return nil
	}

	// 非字符串类型，不能替换
	var value string
	if value = f.String(); len(value) == 0 {
		return nil
	}

	re, err := regexp.Compile(`(?Uis)(href|src)=["']?(.*)["']?.*>`)
	if err != nil {
		err = fmt.Errorf("Field.fixValueURL Compile Error: %v", err)
		log.Errorln(err)
		return err
	}
	f.value = re.ReplaceAllStringFunc(value, func(s string) string {
		if len(s) < 5 {
			return s
		}
		isHref := false
		if s[0:5] == "href=" {
			isHref = true
			s = s[5:]
		} else if s[0:4] == "src=" {
			s = s[4:]
		}
		if s[0] == '"' || s[0] == '\'' {
			s = s[1:]
		}
		if p := strings.Index(s, " "); p > 0 {
			s = s[:p]
		}
		if s[len(s)-1] == '>' {
			s = s[0 : len(s)-1]
		}
		if s[len(s)-1] == '/' {
			s = s[0 : len(s)-1]
		}
		if s[len(s)-1] == '"' || s[len(s)-1] == '\'' {
			s = s[0 : len(s)-1]
		}
		if isHref {
			return "href=\"" + u.FixURL(s) + "\">"
		}
		return "src=\"" + u.FixURL(s) + "\" />"
	})

	f.fixed = true

	return nil
}
