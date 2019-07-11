package url

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/safeie/spider/common/log"
	"github.com/safeie/spider/component/fetcher"
)

// Remote 远程字段
// 远程URL，有一个变量即使当前字段的值，使用 {{.}} 做占位符表示
// -- 如果URL就是该字段的值本身，那么 url={{.}}
type Remote struct {
	url         string           // 远程URL
	pageType    int              // 页面类型
	engine      int              // 抓取引擎
	fetchOption *fetcher.Option  // 获取参数
	logger      log.SimpleLogger // 日志器
}

// NewRemote 创建一个新的远程获取
func NewRemote(logger log.SimpleLogger, pageType int, url string) *Remote {
	t := new(Remote)
	t.logger = logger
	t.pageType = pageType
	t.url = url
	t.fetchOption = fetcher.NewOption("")
	return t
}

// FetchURI 获取字段的远程页面
func (t *Remote) FetchURI(v string) (*URI, error) {
	v = strings.Trim(v, "\"")
	fetch := fetcher.New(t.engine, t.fetchOption)
	uri := strings.Replace(t.url, "{{.}}", url.QueryEscape(v), 1)
	if uri == "" {
		return nil, fmt.Errorf("Field.Remote.Fetch remote url is empty")
	}
	u := NewURI(uri)
	u.PageType = t.pageType

	// 处理值占位符
	if params := t.GetParams(); params != nil {
		for pk, pv := range params {
			if pv == "{{.}}" {
				u.SetParam(pk, url.QueryEscape(v))
			} else {
				u.SetParam(pk, pv)
			}
		}
	}
	res, err := fetch.Fetch(u.URL, u.Req.Params, u.Req.Header)
	if err != nil {
		t.logger.Printf("字段远程页面抓取失败 %s: %v", u.URL, err)
		return nil, fmt.Errorf("Field.Remote.Fetch error: %v", err)
	}
	t.logger.Printf("字段远程页面抓取成功 %s", u.URL)

	u.Code = res.Code
	u.Body = res.Body
	u.Header = res.Header

	return u, nil
}

// EnableJS 是否启用JS渲染
func (t *Remote) EnableJS(ok bool) *Remote {
	if ok {
		t.engine = fetcher.EngineWebKit
		t.SetCharset("UTF-8") // Webkit Phantomjs渲染情况下，自动utf-8
	} else {
		t.engine = fetcher.EngineGoKit
	}
	return t
}

// SetMethod 设置HTTP请求方法
func (t *Remote) SetMethod(v string) *Remote {
	t.fetchOption.SetMethod(v)
	return t
}

// GetMethod 获取HTTP请求方法
func (t *Remote) GetMethod() string {
	return t.fetchOption.GetMethod()
}

// SetHeader 设置HTTP请求头信息
func (t *Remote) SetHeader(key, val string) *Remote {
	t.fetchOption.SetHeader(key, val)
	return t
}

// SetParam 设置HTTP请求附加参数
func (t *Remote) SetParam(key, val string) *Remote {
	t.fetchOption.SetParam(key, val)
	return t
}

// GetParams 获取HTTP请求附加参数
func (t *Remote) GetParams() map[string]string {
	return t.fetchOption.GetParams()
}

// SetCookie 设置Cookie
func (t *Remote) SetCookie(v string) *Remote {
	t.fetchOption.SetCookie(v)
	return t
}

// SetCharset 设置页面编码，默认UTF-8
func (t *Remote) SetCharset(v string) *Remote {
	t.fetchOption.SetCharset(v)
	return t
}

// SetProxy 设置代理服务：类型，地址
func (t *Remote) SetProxy(typ int, addr string) *Remote {
	t.fetchOption.SetProxy(typ, addr)
	return t
}

// SetUserAgent 设置UserAgent
func (t *Remote) SetUserAgent(typ int, ua string) *Remote {
	t.fetchOption.SetUserAgent(typ, ua)
	return t
}

// SetRenderDelay 设置JS渲染时间，单位是 秒，默认不等待，如果发现页面JS没有执行完整，可以加大该时间
func (t *Remote) SetRenderDelay(v int) *Remote {
	t.fetchOption.SetRenderDelay(v)
	return t
}

// SetTimeout 设置抓取超时，默认 1秒
func (t *Remote) SetTimeout(v int) *Remote {
	t.fetchOption.SetTimeout(v)
	return t
}
