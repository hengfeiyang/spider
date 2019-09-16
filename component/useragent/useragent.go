package useragent

const (
	Custom  = iota // 自定义
	Rand           // 随机
	Common         // 普通，通用
	PC             // 电脑
	Mobile         // 手机
	IOS            // iOS
	IPhone         // iPhone
	IPad           // iPad
	MacOS          // macOS
	Android        // Android
	Wechat         // Wechat
	QQ             // QQ
	Baidu          // spider, Baidu
	Google         // spider, Google
	Bing           // spider, Bing
	Sogou          // spider, Sogou
	Qihu           // spider, Qihu
	Yahoo          // spider, Yahoo
	last           // last
)

var name []string
var ua []string

func init() {
	name = make([]string, last+1)
	name[Custom] = "自定义"
	name[Common] = "通用"
	name[PC] = "电脑"
	name[Mobile] = "手机"
	name[IOS] = "iOS"
	name[IPhone] = "iPhone"
	name[IPad] = "iPad"
	name[MacOS] = "macOS"
	name[Android] = "Android"
	name[Wechat] = "微信"
	name[QQ] = "QQ"
	name[Baidu] = "百度蜘蛛"
	name[Google] = "谷歌蜘蛛"
	name[Bing] = "必应蜘蛛"
	name[Sogou] = "搜狗蜘蛛"
	name[Qihu] = "360蜘蛛"
	name[Yahoo] = "雅虎蜘蛛"

	ua = make([]string, last+1)
	ua[Common] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/37.0.963.56 Safari/535.11"
	ua[PC] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/37.0.963.56 Safari/535.11"
	ua[Mobile] = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	ua[IOS] = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	ua[IPhone] = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	ua[IPad] = "Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	ua[MacOS] = "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	ua[Android] = "Mozilla/5.0 (Linux; U; Android 5.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"
	ua[Wechat] = "Mozilla/5.0 (Linux; U; Android 5.4.2; zh-cn; 2014011 Build/HM2014011) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025489 Mobile Safari/533.1 MicroMessenger/6.3.15.49_r8aff805.760 NetType/WIFI Language/zh_CN"
	ua[QQ] = "Mozilla/5.0 (Linux; Android 5.0; F103 Build/LRX21M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/37.0.0.0 Mobile Safari/537.36 V1_AND_SQ_6.2.3_336_YYB_D QQ/6.2.3.2700 NetType/WIFI WebP/0.4.1 Pixel/720"
	ua[Baidu] = "Baiduspider+(+http://www.baidu.com/search/spider.htm)"
	ua[Google] = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	ua[Bing] = "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"
	ua[Sogou] = "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)"
	ua[Qihu] = "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0); 360Spider"
	ua[Yahoo] = "Mozilla/5.0 (compatible; Yahoo! Slurp China; http://misc.yahoo.com.cn/help.html)"
}

// Name 获取UA类型名称
func Name(t int) string {
	s := name[t]
	if s == "" {
		return name[Custom]
	}
	return s
}

// Get 获取UserAgent
func Get(t int) string {
	s := ua[t]
	if s == "" {
		return ua[Common]
	}
	return s
}
