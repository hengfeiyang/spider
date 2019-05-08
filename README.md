# spider
spider is some components for crawl webpage, it is useful.

## common

util functions for spider

## config

runtime config file

## components

### fetcher

* gokit: use go http fetch data
* webkit: use webkit(phantomjs) fetch data, this can parse javascript in webpage

### parser

* htmldom: use goquery parse html document, support dom xpath
* jsonpath: use json parse json data
* regexp: use regexp parse data
* substring: use split and substr to parse data

### proxy

* kxdaili: provider proxy service for http request use kx100.com
* be more

### task

provider a full task crawl project, start with url and end with data

### url

provider url filter, fixpath, fetch content, parse content and more functions

### useragent

* Common         // 普通，通用
* PC             // 电脑
* Mobile         // 手机
* IOS            // iOS
* IPhone         // iPhone
* IPad           // iPad
* MacOS          // macOS
* Android        // Android
* Wechat         // Wechat
* QQ             // QQ
* Baidu          // spider, Baidu
* Google         // spider, Google
* Bing           // spider, Bing
* Sogou          // spider, Sogou
* Qihu           // spider, Qihu
* Yahoo          // spider, Yahoo
