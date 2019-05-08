package fetcher

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/safeie/spider/component/proxy"
	"github.com/safeie/spider/component/useragent"
	iconv "gopkg.in/iconv.v1"
)

const (
	// EngineGoKit Go直接抓取
	EngineGoKit = iota
	// EngineWebKit Webkit渲染
	EngineWebKit
)

// Fetcher fetch interface
type Fetcher interface {
	// Fetch 执行请求，返回 body,header, cookie, error
	Fetch(url string) (*Response, error)
}

// Option 抓取器的配置参数
type Option struct {
	rootDir       string            // 目录，执行目录，phantomjs和脚本将从目录中获取
	method        string            // HTTP请求方法
	headers       map[string]string // Header头设置
	params        map[string]string // HTTP请求附加字段
	cookieLock    sync.RWMutex      // Cookie锁
	_cookie       string            // Cookie信息，禁止直接使用，所以带了个下划线
	charset       string            // 网站页面编码
	userAgentType int               // UserAgent类型
	userAgent     string            // 自定义UserAgent
	proxyType     int               // 代理类型设置，不使用代理，自定义代理，启用国内代理，启用国外代理
	proxyAddr     string            // 代理服务器地址
	renderDelay   int               // 渲染等待，单位 毫秒，用于js渲染时获取内容前的等待，确保渲染完成
	timeout       int               // 抓取超时，单位 秒
}

// NewOption 创新新的抓取配置
func NewOption(rootDir string) *Option {
	if rootDir == "" {
		_, filename, _, _ := runtime.Caller(1)
		rootDir = path.Join(path.Dir(filename), "../")
	}

	t := new(Option)
	t.rootDir = rootDir
	t.method = "GET"
	t.params = make(map[string]string)
	t.headers = make(map[string]string)
	t.charset = "UTF-8"
	t.renderDelay = 100 // 0.1秒
	t.timeout = 5       // 默认5秒超时
	return t
}

// SetMethod 设置HTTP请求方法
func (t *Option) SetMethod(v string) {
	if v != "POST" {
		v = "GET"
	}
	t.method = v
}

// SetHeader 设置HTTP请求头信息
func (t *Option) SetHeader(key, val string) {
	t.headers[key] = val
}

// SetParam 设置HTTP请求附件参数
func (t *Option) SetParam(key, val string) {
	t.params[key] = val
}

// GetCookie 获取cookie信息，因为带锁，所以要通过两个方法获取
func (t *Option) GetCookie() string {
	t.cookieLock.RLock()
	c := t._cookie
	t.cookieLock.RUnlock()
	return c
}

// SetCookie 设置Cookie !唯一一个可以在运行中修改的参数
func (t *Option) SetCookie(v string) {
	t.cookieLock.Lock()
	t._cookie = v
	t.cookieLock.Unlock()
}

// GetCharset 获取页面编码
func (t *Option) GetCharset() string {
	return t.charset
}

// SetCharset 设置页面编码，默认UTF-8
func (t *Option) SetCharset(v string) {
	t.charset = v
}

// GetProxy 获取代理服务器
func (t *Option) GetProxy() (int, string) {
	return t.proxyType, t.proxyAddr
}

// SetProxy 设置代理服务：类型，地址
func (t *Option) SetProxy(typ int, addr string) {
	t.proxyType = typ
	if typ == proxy.TypeCustom && addr != "" {
		t.proxyAddr = addr
	}
}

// SetUserAgent 设置UserAgent
func (t *Option) SetUserAgent(typ int, ua string) {
	t.userAgentType = typ
	if typ == useragent.Custom && ua != "" {
		t.userAgent = ua
	}
}

// SetRenderDelay 设置JS渲染时间，单位是 秒，默认不等待，如果发现页面JS没有执行完整，可以加大该时间
func (t *Option) SetRenderDelay(v int) {
	t.renderDelay = v
}

// SetTimeout 设置抓取超时，默认 1秒
func (t *Option) SetTimeout(v int) {
	if v == 0 {
		v = 1
	}
	t.timeout = v
}

// GetProxyAddr 获取代理IP地址
func GetProxyAddr(option *Option) string {
	// 不使用代理
	if option.proxyType == proxy.TypeNone {
		return ""
	}
	// 检测自定义代理
	if option.proxyAddr != "" {
		return option.proxyAddr
	}
	// 获取代理
	p := proxy.Get(option.proxyType)
	if p == "" {
		return ""
	}
	if !strings.Contains(p, "://") {
		p = "http://" + p
	}
	return p
}

// GetUserAgent 获取浏览器代理头
func GetUserAgent(option *Option) string {
	// 自定义
	if option.userAgentType == useragent.Custom && option.userAgent != "" {
		return option.userAgent
	}
	// 系统内置
	return useragent.Get(option.userAgentType)
}

// Response 抓取器的返回结果
type Response struct {
	Code   int         // 返回的状态码
	Cookie string      // 返回的会话信息
	Body   []byte      // 返回的内容
	Header http.Header // 返回的头部
}

// New 创建一个抓取器
func New(kit int, option *Option) Fetcher {
	if kit == EngineWebKit {
		return &Webkit{option}
	}
	return &Gokit{option}
}

// ConvertCharset convert str from charset to UTF-8
func ConvertCharset(res *Response, charset string) ([]byte, error) {
	if charset == "" {
		charset = GetCharset(res)
	}
	if charset == "UTF-8" {
		return res.Body, nil
	}
	cd, err := iconv.Open("UTF-8", charset)
	if err != nil {
		return nil, err
	}
	newBody := cd.ConvString(string(res.Body))
	return []byte(newBody), nil
}

// GetCharset return the url charset, when get empty log error, and return UTF-8
func GetCharset(res *Response) string {
	if val, ok := res.Header["Content-Type"]; ok {
		var _val string
		var _tmp []string
		for _, v := range val {
			if strings.Index(v, ";") > -1 {
				_tmp = strings.Split(v, ";")
			} else {
				_tmp = []string{v}
			}
			for i := range _tmp {
				_val = strings.TrimSpace(strings.ToUpper(_tmp[i]))
				if strings.HasPrefix(_val, "CHARSET=") {
					_val = strings.Trim(_val[8:], "\"")
					return _val
				}
			}
		}
	}
	charset := GetCharsetFromHTML(res.Body)
	if charset == "" {
		return "UTF-8"
	}
	return charset
}

// GetCharsetFromHTML from html find the charset
func GetCharsetFromHTML(body []byte) string {
	re, err := regexp.Compile("(?is)<meta .*?charset=['\"]?([0-9a-zA-Z-]+)(.*?)>")
	if err != nil {
		return ""
	}
	matches := re.FindStringSubmatch(string(body))
	if matches == nil || len(matches) < 2 {
		return ""
	}
	charset := matches[1]
	return strings.TrimSpace(strings.ToUpper(charset))
}

// FixURL 修复相对路径
func FixURL(base, href string) string {
	if href == "" {
		return ""
	}

	if href[0] == '.' && len(href) > 1 && href[1] == '/' {
		href = href[2:]
	}

	// 过滤锚点
	if href[0] == '#' {
		return ""
	}
	if pos := strings.Index(href, "#"); pos > 0 {
		href = href[0:pos]
	}
	// 过滤JS连接
	if strings.Contains(href, "javascript:") {
		return ""
	}
	// 过滤mailto
	if strings.Contains(href, "mailto:") {
		return ""
	}
	// 完整的URL
	if strings.Contains(href, "://") {
		return href
	}

	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	// 补齐URL
	if href[0] == '/' {
		if len(href) > 1 && href[1] == '/' {
			// 补齐协议，现在流行这种URL  //u.domain.com
			return u.Scheme + ":" + href
		}
		return u.Scheme + "://" + u.Host + href
	}
	return u.Scheme + "://" + u.Host + u.Path + "/" + href
}
