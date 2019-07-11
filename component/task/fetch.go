package task

import (
	"sync"

	"github.com/safeie/spider/component/fetcher"
	"github.com/safeie/spider/component/url"
)

const (
	// MethodGET GET方法获取
	MethodGET = "GET"
	// MethodPOST POST方法获取
	MethodPOST = "POST"
)

// CheckRepeatFunc 检测重复方法
type CheckRepeatFunc func(uri *url.URI) bool

// BeforeFetchFunc 获取前置方法
type BeforeFetchFunc func(uri *url.URI)

// AfterFetchFunc 获取后置方法
type AfterFetchFunc func(uri *url.URI)

// AntiSpiderFunc 判断是否触发了反作弊策略
type AntiSpiderFunc func(uri *url.URI) bool

// FetcherPool 抓取器池
type FetcherPool struct {
	kit   int
	task  *Task
	store chan fetcher.Fetcher
	mu    sync.Mutex
}

// NewFetcherPool 创建一个新的抓取器池子
func NewFetcherPool(initCap, maxCap int, kit int, task *Task) *FetcherPool {
	if maxCap == 0 || initCap > maxCap {
		maxCap = initCap * 2
	}
	p := new(FetcherPool)
	p.kit = kit
	p.task = task
	p.store = make(chan fetcher.Fetcher, maxCap)
	for i := 0; i < initCap; i++ {
		p.store <- p.create()
	}
	return p
}

func (p *FetcherPool) create() fetcher.Fetcher {
	return fetcher.New(p.kit, p.task.setting.fetchOption)
}

// Get 取一个
func (p *FetcherPool) Get() fetcher.Fetcher {
	for {
		select {
		case v := <-p.store:
			return v
		default:
			return p.create()
		}
	}
}

// Put 还回
func (p *FetcherPool) Put(f fetcher.Fetcher) {
	select {
	case p.store <- f:
		return
	default:
		return
	}
}
