package proxy

import (
	"github.com/safeie/spider/component/proxy/base"
	"github.com/safeie/spider/component/proxy/kxdaili"
)

const (
	TypeNone    = iota // 不使用代理
	TypeCustom         // 自定义代理
	TypeInside         // 国内代理
	TypeOutside        // 国外代理
	last
)

var name []string
var _defaultProxy *proxy

func init() {
	name = make([]string, last+1)
	name[TypeNone] = "不使用"
	name[TypeCustom] = "自定义"
	name[TypeInside] = "国内代理"
	name[TypeOutside] = "国外代理"

	_defaultProxy = new(proxy)
	_defaultProxy.provider = kxdaili.New()
}

// Name 获取代理类型名称
func Name(t int) string {
	s := name[t]
	if s == "" {
		return name[TypeNone]
	}
	return s
}

// Get 获取一个指定类型的proxy地址
func Get(t int) string {
	return _defaultProxy.get(t)
}

// proxy 代理管理器
type proxy struct {
	provider base.Proxyer
}

// get 获取一个指定类型的proxy地址
func (p *proxy) get(t int) string {
	switch t {
	case TypeNone:
		return ""
	case TypeCustom:
		return ""
	case TypeInside:
	case TypeOutside:
	default:
	}
	return p.provider.Get()
}
