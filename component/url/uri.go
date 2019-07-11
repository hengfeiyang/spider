package url

import (
	"net/http"
	"net/url"

	"github.com/safeie/spider/component/fetcher"
	"github.com/safeie/spider/component/parser"
)

const (
	// PageTypeHTML 页面类型，HTML网页
	PageTypeHTML = iota
	// PageTypeJSON 页面类型，JSON数据
	PageTypeJSON
	// PageTypeText 页面类型，纯文本
	PageTypeText
)

// URI 是URL的组成单元，不直接使用string的原因是可以附加数据
type URI struct {
	URL       string                 // URL
	parsedURL *url.URL               // 标准的URL解析
	PageType  int                    // 页面类型
	Code      int                    // 请求的响应码
	Header    http.Header            // 请求的响应头
	Body      []byte                 // 请求的响应体
	Fetched   bool                   // 是否抓取过
	fields    []*Field               // 字段
	attach    map[string]interface{} // 附加数据
	Req       struct {               // 请求参数
		Header map[string]string
		Params map[string]string
	}
	Parser struct { // 解析器
		JSON      *parser.JSONPath
		Regexp    *parser.Regexp
		Substring *parser.Substring
		HTMLDom   *parser.HTMLDom
	}
}

// NewURI 创建一个新的URI结构
func NewURI(uri string) *URI {
	u := new(URI)
	u.URL = uri
	u.parsedURL, _ = url.Parse(uri)
	u.attach = make(map[string]interface{})
	u.Req.Header = make(map[string]string)
	u.Req.Params = make(map[string]string)
	return u
}

// Copy 复制一个URI结构，不复制 Body, Parser, fields
func (u *URI) Copy() *URI {
	n := new(URI)
	n.URL = u.URL
	n.parsedURL = u.parsedURL
	n.PageType = u.PageType
	n.Code = u.Code
	n.Header = u.Header
	n.Fetched = u.Fetched
	n.attach = u.attach
	n.Req.Header = u.Req.Header
	n.Req.Params = u.Req.Params
	return n
}

// ResetBody 重设以复用该实例
func (u *URI) ResetBody(body []byte) {
	u.Body = u.Body[:0]
	u.Body = append(u.Body, body...)
	u.Parser.JSON = nil
	u.Parser.Regexp = nil
	u.Parser.Substring = nil
	u.Parser.HTMLDom = nil
}

// Set 设置一个附加属性
func (u *URI) Set(key string, val interface{}) {
	u.attach[key] = val
}

// Get 获取保存的附加属性
func (u *URI) Get(key string) interface{} {
	return u.attach[key]
}

// Gets 返回附加的属性列表
func (u *URI) Gets() map[string]interface{} {
	return u.attach
}

// SetHeader 设置URL请求时附属的Header头
func (u *URI) SetHeader(key, val string) {
	u.Req.Header[key] = val
}

// SetParam 设置URL请求时附属的请求参数
func (u *URI) SetParam(key, val string) {
	u.Req.Params[key] = val
}

// FetchURLs 获取内容中的URL列表
func (u *URI) FetchURLs() []string {
	if u.PageType != PageTypeHTML {
		return nil
	}

	if len(u.Body) == 0 {
		return nil
	}

	if u.Parser.HTMLDom == nil {
		var err error
		u.Parser.HTMLDom, err = parser.NewHTMLDom(u.Body)
		if err != nil {
			return nil
		}
	}
	var urls []string
	as := u.Parser.HTMLDom.DomFind("a")
	for i := 0; i < as.Size(); i++ {
		if href := u.Parser.HTMLDom.NodeAttr(as.Get(i), "href"); href != "" {
			urls = append(urls, u.FixURL(href))
		}
	}

	return urls
}

// AddFields 添加字段
func (u *URI) AddFields(ff ...*Field) {
	u.fields = append(u.fields, ff...)
}

// ExportFields 导出该URL保存的字段内容
func (u *URI) ExportFields() map[string]interface{} {
	store := make(map[string]interface{})
	for _, f := range u.fields {
		store[f.Name] = f.Value()
	}
	return store
}

// FixURL 过滤非本站URL和修复相对路径
func (u *URI) FixURL(href string) string {
	return fetcher.FixURL(u.URL, href)
}
