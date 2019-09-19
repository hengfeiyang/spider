package url

import (
	"sync"
)

// URL URL组件
/*
 * URL包含两种URL：
 * initURL 即入口URL，会无条件提取整个页面上的URL添加到待处理列表
 * ruleURL，处理过程中产生的 规则URL，这些URL要经过规则处理，不符合规则的被丢弃
 */
type URL struct {
	inited   bool            // 判断是否初试化的标志位
	initfunc URLinitFunc     // 初始化函数
	initURLs []string        // 入口URL，无条件抓取所有的url
	ruleURLs chan *URI       // 待处理的URL列表，所有的URL将经过规则处理
	crawled  map[string]bool // 存储抓取过的URL，用于判断是否还要抓取
	rw       sync.Mutex
}

// URLinitFunc URL初试化函数
type URLinitFunc func() []string

// NewURL 创建一个新的URL实例
func NewURL() *URL {
	t := new(URL)
	t.initURLs = make([]string, 0)
	t.ruleURLs = make(chan *URI, 100000)
	t.crawled = make(map[string]bool)
	return t
}

// Reset 复用
func (t *URL) Reset() *URL {
	t.rw.Lock()
	t.crawled = make(map[string]bool)
	t.rw.Unlock()
	return t
}

// SetInitFunc 设置URL初始化函数，用于自定义初始化URL列表，限制最大1000个初始URL
func (t *URL) SetInitFunc(f URLinitFunc) {
	t.initfunc = f
}

// GetInitURLs 获取入口URL
func (t *URL) GetInitURLs() []string {
	return t.initURLs
}

// PushURL 插入一个ruleURL
// 判断重复，防止进入无限循环，采集过就不再入队列
func (t *URL) PushURL(url string) {
	t.rw.Lock()
	if t.crawled[url] {
		t.rw.Unlock()
		return
	}
	t.crawled[url] = true
	t.rw.Unlock()
	t.ruleURLs <- NewURI(url)
}

// Push 插入一个ruleURL URI结构
// 判断重复，防止进入无限循环，采集过就不再入队列
func (t *URL) Push(uri *URI) {
	t.rw.Lock()
	if t.crawled[uri.URL] {
		t.rw.Unlock()
		return
	}
	t.crawled[uri.URL] = true
	t.rw.Unlock()
	t.ruleURLs <- uri
}

// Pop 获取一个ruleURL，并从存储中删除
func (t *URL) Pop() *URI {
	return <-t.ruleURLs
}

// Len 返回当前ruleURL数量
func (t *URL) Len() int {
	return len(t.ruleURLs)
}

// Initialize 初始化URL管理器
func (t *URL) Initialize() {
	if t.inited {
		return
	}
	if t.initfunc != nil {
		if urls := t.initfunc(); urls != nil {
			for i := range urls {
				t.initURLs = append(t.initURLs, urls[i])
			}
		}
	}
	t.inited = true
}
