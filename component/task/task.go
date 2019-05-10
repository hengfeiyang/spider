// Package task 提供任务管理
package task

import (
	"sync"
	"time"

	"github.com/safeie/spider/common/log"
	"github.com/safeie/spider/component/fetcher"
	"github.com/safeie/spider/component/url"
)

const (
	stopEmpty        int = iota
	stopNotSaveQueue     // 停止不保存队列
	stopAndSaveQueue     // 停止并保存队列
)

// Task 任务
type Task struct {
	id          string        // 任务编号
	name        string        // 任务名称
	domain      string        // 该任务匹配的域名
	configDir   string        // 配置目录，用于获取配置脚本
	logger      log.Logger    // 日志器
	url         *url.URL      // URL管理器
	fetcherPool *FetcherPool  // 抓取器
	rule        []*Rule       // 采集规则
	setting     *taskSetting  // 设置项
	lockrunning sync.RWMutex  // 任务运行锁
	running     bool          // 任务是否在运行中
	routineNum  int           // 协程数量，默认10个
	chanLink    chan struct{} // 控制抓取的协程通道
	shutdown    chan int      // 关闭，根据传递的信号决定是否保存队列
	urlStore    *sync.Map     // url存储器，用来判断是否已经采集过
}

// taskSetting 任务设置项目，可外部设置的内容
type taskSetting struct {
	fetchOption     *fetcher.Option // 抓取，配置
	engine          int             // 抓取，抓取引擎
	interval        int             // 执行间隔，单位 毫秒，用于限制采集频率
	autoSession     bool            // 是否自动记录会话
	errorContinue   bool            // 出错后，是否继续下一个URL
	retryTimes      int             // 出错，重试次数
	prepareFunc     PrepareFunc     // 任务，预处理钩子函数
	antiSpiderFunc  AntiSpiderFunc  // 抓取，反作弊函数
	checkRepeatFunc CheckRepeatFunc // 抓取，重复检测钩子函数
	beforeFetchFunc BeforeFetchFunc // 抓取，前置钩子函数
	afterFetchFunc  AfterFetchFunc  // 抓取，后置钩子函数
	beforeQuitFunc  BeforeQuitFunc  // 停止，前置钩子函数
}

// PrepareFunc 任务预处理函数
type PrepareFunc func(*Task)

// New 创建新的任务，参数为：任务编号，名称，域名，脚本目录
//
// configDir 可为空，默认值是代码目录中的 config 目录
func New(taskID, taskName, domain, configDir string) *Task {
	t := new(Task)
	t.id = taskID
	t.name = taskName
	t.domain = domain
	t.configDir = configDir
	t.logger = log.DefaultLogger
	t.url = url.NewURL()
	t.rule = make([]*Rule, 0, 2)
	t.routineNum = 10
	t.shutdown = make(chan int)

	t.setting = new(taskSetting)
	t.setting.interval = 100 // 0.1秒
	t.setting.autoSession = false
	t.setting.errorContinue = false
	t.setting.retryTimes = 3

	// 抓取设置
	t.setting.fetchOption = fetcher.NewOption(t.configDir)

	// 初始化存储
	t.urlStore = new(sync.Map)

	return t
}

// Reset 复用
func (t *Task) Reset() *Task {
	t.url.Reset()
	return t
}

// ID 返回任务的编号
func (t *Task) ID() string {
	return t.id
}

// Name 返回任务的名称
func (t *Task) Name() string {
	return t.name
}

// Domain 返回任务的域名
func (t *Task) Domain() string {
	return t.domain
}

// Log 输出日志调试
func (t *Task) Log(v ...interface{}) {
	t.logger.Print(v...)
}

// Logf 输出日志调试
func (t *Task) Logf(format string, v ...interface{}) {
	t.logger.Printf(format, v...)
}

// Logln 输出日志调试
func (t *Task) Logln(v ...interface{}) {
	t.logger.Println(v...)
}

// SetRoutineNum 设置任务的协程数量，默认 10个
func (t *Task) SetRoutineNum(v int) *Task {
	if v == 0 {
		v = 10
	}
	t.routineNum = v
	return t
}

// EnableJS 是否启用JS渲染
func (t *Task) EnableJS(ok bool) *Task {
	if ok {
		t.setting.engine = fetcher.EngineWebKit
		t.SetCharset("UTF-8") // Webkit Phantomjs渲染情况下，自动utf-8
	} else {
		t.setting.engine = fetcher.EngineGoKit
	}
	return t
}

// SetInterval 设置采集间隔，单位 微妙，默认 100微秒
func (t *Task) SetInterval(v int) *Task {
	t.setting.interval = v
	return t
}

// SetAutoSession 设置是否自动记录会话
// 比如，雪球网，必须先访问一下HTML页面记录下会话才可以继续请求JSON数据
// 比如，豆瓣网，根据cookie会话统计访问频次，不能记录cookie
func (t *Task) SetAutoSession(v bool) *Task {
	t.setting.autoSession = v
	return t
}

// SetErrorContinue 设置遇到错误是，是否继续，默认 遇到错误即终止运行，并保存队列中的数据
func (t *Task) SetErrorContinue(v bool) *Task {
	t.setting.errorContinue = v
	return t
}

// SetMethod 设置HTTP请求方法
func (t *Task) SetMethod(v string) *Task {
	t.setting.fetchOption.SetMethod(v)
	return t
}

// SetHeader 设置HTTP请求头信息
func (t *Task) SetHeader(key, val string) *Task {
	t.setting.fetchOption.SetHeader(key, val)
	return t
}

// SetParam 设置HTTP请求附件参数
func (t *Task) SetParam(key, val string) *Task {
	t.setting.fetchOption.SetParam(key, val)
	return t
}

// SetCookie 设置Cookie
func (t *Task) SetCookie(v string) *Task {
	t.setting.fetchOption.SetCookie(v)
	return t
}

// SetCharset 设置页面编码，默认UTF-8
func (t *Task) SetCharset(v string) *Task {
	t.setting.fetchOption.SetCharset(v)
	return t
}

// SetProxy 设置代理服务：类型，地址
func (t *Task) SetProxy(typ int, addr string) *Task {
	t.setting.fetchOption.SetProxy(typ, addr)
	return t
}

// SetUserAgent 设置UserAgent：类型，自定义
func (t *Task) SetUserAgent(typ int, ua string) *Task {
	t.setting.fetchOption.SetUserAgent(typ, ua)
	return t
}

// SetRenderDelay 设置JS渲染时间，单位是 秒，默认不等待，如果发现页面JS没有执行完整，可以加大该时间
func (t *Task) SetRenderDelay(v int) *Task {
	t.setting.fetchOption.SetRenderDelay(v)
	return t
}

// SetTimeout 设置抓取超时，默认 1秒
func (t *Task) SetTimeout(v int) *Task {
	t.setting.fetchOption.SetTimeout(v)
	return t
}

// SetPrepareFunc 设置开始执行前的预处理函数
func (t *Task) SetPrepareFunc(p PrepareFunc) *Task {
	t.setting.prepareFunc = p
	return t
}

// SetFetchFunc 设置抓取的钩子函数
func (t *Task) SetFetchFunc(c CheckRepeatFunc, b BeforeFetchFunc, a AfterFetchFunc) *Task {
	t.setting.checkRepeatFunc = c
	t.setting.beforeFetchFunc = b
	t.setting.afterFetchFunc = a
	return t
}

// SetBeforeQuitFunc 设置停止前置钩子函数
func (t *Task) SetBeforeQuitFunc(f BeforeQuitFunc) *Task {
	t.setting.beforeQuitFunc = f
	return t
}

// SetAntiSpiderFunc 设置抓取的反作弊检测函数
func (t *Task) SetAntiSpiderFunc(f AntiSpiderFunc) *Task {
	t.setting.antiSpiderFunc = f
	return t
}

// SetURLinitFunc 设置URL初始化函数
func (t *Task) SetURLinitFunc(f url.URLinitFunc) *Task {
	t.url.SetInitFunc(f)
	return t
}

// Rule 设置一个执行规则并返回
// 同一个URL应该仅能匹配到一个规则，执行时，仅处理匹配到的第一个规则
func (t *Task) Rule(s string) *Rule {
	for i := range t.rule {
		if t.rule[i].rule == s {
			return t.rule[i]
		}
	}

	r := newRule(s, t)
	t.rule = append(t.rule, r)
	return r
}

// FetchURL 获取单个网页，给外部调用
func (t *Task) FetchURL(uri string) (u *url.URI, cookie string, err error) {
	u = url.NewURI(uri)
	cookie, err = t.FetchURI(u, t.fetcherPool)
	return u, cookie, err
}

// FetchURI 执行一个URI的内容获取，内部使用，如果出错，至少尝试N次
func (t *Task) FetchURI(u *url.URI, fetcherPool *FetcherPool) (string, error) {
	var err error
	var res *fetcher.Response
	f := fetcherPool.Get()
	for i := 0; i < t.setting.retryTimes; i++ {
		res, err = f.Fetch(u.URL)
		if res != nil && (res.Code == 404 || res.Code == 403) {
			break
		}
		if err == nil {
			break
		}
	}
	fetcherPool.Put(f)
	if err != nil {
		t.Logf("页面抓取失败 %s: %v", u.URL, err)
	} else {
		t.Logf("页面抓取成功 %s", u.URL)
	}

	if err != nil {
		return "", err
	}

	// 自动记录会话
	if t.setting.autoSession && res.Cookie != "" && t.setting.fetchOption.GetCookie() == "" {
		t.SetCookie(res.Cookie)
	}

	// 设置抓取过
	u.Code = res.Code
	u.Body = res.Body
	u.Header = res.Header
	u.Fetched = true

	return res.Cookie, nil
}

// PushURL 向任务过程中添加URL字符串，接受规则校验，按规则处理，不符合规则的将丢弃
func (t *Task) PushURL(urls ...string) *Task {
	for i := range urls {
		if rule := t.matchRule(urls[i]); rule != nil {
			t.url.PushURL(urls[i])
		}
	}
	return t
}

// PushURI 向任务过程中添加URI结构，接受规则校验，按规则处理，不符合规则的将丢弃
func (t *Task) PushURI(uris ...*url.URI) *Task {
	for i := range uris {
		if rule := t.matchRule(uris[i].URL); rule != nil {
			t.url.Push(uris[i])
		}
	}
	return t
}

// matchRule 是否匹配到可用规则
func (t *Task) matchRule(url string) []*Rule {
	var rs []*Rule
	for i := range t.rule {
		if ok := t.rule[i].Match(url); ok {
			rs = append(rs, t.rule[i])
		}
	}
	return rs
}

// runRule 使用匹配的规则执行绑定的方法
func (t *Task) runRule(uri *url.URI, fetcherPool *FetcherPool) {
	var err error
	if rules := t.matchRule(uri.URL); rules != nil {
		for _, r := range rules {
			if err = r.Run(uri, fetcherPool); err != nil {
				if !t.setting.errorContinue {
					t.Logf("任务遇到错误即将终止: %v", err)
					go func() {
						t.Stop()
					}()
					break
				}
			}
		}
	}
}

// Run 开始执行
func (t *Task) Run() error {
	// 运行中，禁止修改
	if t.isRunning() {
		return Errorf("Task is running...")
	}
	t.setRunning(true)
	defer func() {
		t.setRunning(false)
	}()

	// 初始化URL控制器
	t.url.Initialize()
	// 初始化抓取器
	t.fetcherPool = NewFetcherPool(t.routineNum, 0, t.setting.engine, t)
	// 控制抓取的协程数
	t.chanLink = make(chan struct{}, t.routineNum)
	// 运行prepare函数
	if t.setting.prepareFunc != nil {
		t.setting.prepareFunc(t)
	}

	// 开启种子协程
	for _, u := range t.url.GetInitURLs() {
		go func(u string) {
			t.PushURL(u)
		}(u)
	}

	// 开启主进程
	shutdown := 0
	ticker := time.NewTicker(time.Millisecond * time.Duration(t.setting.interval))
	for {
		select {
		case ch := <-t.shutdown:
			shutdown = ch
		case <-ticker.C:
			// get a ticket
		}
		if shutdown > 0 {
			// 停止服务
			if shutdown == stopAndSaveQueue {
				if err := t.save(); err != nil {
					log.Errorln(err)
				}
			}
			ticker.Stop()
			break
		}

		// 没有任务了，退出
		if t.url.Len() == 0 && len(t.chanLink) == 0 {
			break
		}

		// 已经没有URL可执行，等待其他执行协程结束
		if t.url.Len() == 0 {
			continue
		}

		u := t.url.Pop()
		if u == nil || u.URL == "" || u.Fetched {
			continue
		}

		// 抓取内容
		t.chanLink <- struct{}{}
		go func(u *url.URI, ch chan struct{}, fetcherPool *FetcherPool) {
			t.runRule(u, fetcherPool)
			<-ch
		}(u, t.chanLink, t.fetcherPool)
	}

	// 执行结束
	return nil
}

// Stop 停止任务，保存队列
func (t *Task) Stop() {
	if t.isRunning() {
		t.shutdown <- stopAndSaveQueue
	}
}

// Close 关闭任务，直接关闭不保存队列
func (t *Task) Close() {
	if t.isRunning() {
		t.shutdown <- stopNotSaveQueue
	}
}

// isRunning 判断任务是否在运行
func (t *Task) isRunning() bool {
	t.lockrunning.RLock()
	running := t.running
	t.lockrunning.RUnlock()
	return running
}

// setRunning 设置任务运行状态
func (t *Task) setRunning(v bool) {
	t.lockrunning.Lock()
	t.running = v
	t.lockrunning.Unlock()
}

// save 保存任务状态，保存未执行完的URL为临时队列
func (t *Task) save() error {
	queue := make([]string, 0, t.url.Len())
	for {
		time.Sleep(time.Millisecond)

		// 没有任务了，退出
		if t.url.Len() == 0 && len(t.chanLink) == 0 {
			break
		}
		// 已经没有URL可执行，等待其他执行协程结束
		if t.url.Len() == 0 {
			continue
		}
		u := t.url.Pop()
		queue = append(queue, u.URL)
	}
	if len(queue) == 0 {
		return nil
	}

	if t.setting.beforeQuitFunc != nil {
		t.setting.beforeQuitFunc(t.id, queue)
	}

	return nil
}

// logURL 记录url用于重复判断 如果已存在返回 true
func (t *Task) logURL(url, hash string) bool {
	v, ok := t.urlStore.Load(url)
	if ok {
		if hash == v.(string) {
			return true
		}
	}

	t.urlStore.Store(url, hash)
	return false
}
