package kxdaili

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/safeie/spider/common/util"
	"github.com/safeie/spider/component/proxy/base"
)

// KXDaili 开心代理
type KXDaili struct {
	gateway   string
	store     []string
	storeLock sync.RWMutex
	hitTime   time.Time
}

// New 创建一个代理接口
func New() base.Proxyer {
	p := new(KXDaili)
	p.gateway = "http://api.kxdaili.com/?api=201612031700559680&gb=2&ct=500&https=%E6%94%AF%E6%8C%81"
	p.store = make([]string, 0)
	p.fetch()
	return p
}

// Get 获取一个代理IP
// 每3分钟，更新一次IP列表
func (p *KXDaili) Get() string {
	p.storeLock.RLock()
	if p.hitTime.Before(time.Now().Add(time.Second * -180)) {
		go func() {
			p.fetch()
		}()
		p.hitTime = time.Now()
	}

	var v string
	if len(p.store) > 0 {
		n := util.RandInt(len(p.store))
		v = p.store[n]
	}
	p.storeLock.RUnlock()
	return v
}

// fetch 获取IP地址
func (p *KXDaili) fetch() error {
	resp, err := http.Get(p.gateway)
	if err != nil {
		return err
	}
	store := make([]string, 0)
	r := bufio.NewReader(resp.Body)
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}
		store = append(store, string(l))
	}
	resp.Body.Close()
	if len(store) > 0 {
		// 判断是否出错了
		if strings.Contains(store[0], "<bad>") {
			return fmt.Errorf(store[0])
		}
		// 存储
		p.storeLock.Lock()
		p.store = store
		p.storeLock.Unlock()
	}

	return nil
}
