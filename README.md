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

## task flow

回调函数总结：
    task级别
      1. task.PrepareFunc         一次性       传递 TASK，用于初始化, 设置Cookie啥的
      2. task.URLinitFunc         一次性       接收 初始化URLs, 可以批量设置入口URL
      3. task.AntiSpiderFunc      每个链接      传递 URI, 判断一个url是否被反采集拦截
      4. task.BeforeFetchFunc     每个链接      传递 URI, 在获取内容前可以干点啥
      5. task.CheckRepeatFunc     每个链接      传递 URI, 在获取内容前检测是否重复，可用于检测是否是采集过的内容
      6. task.AfterFetchFunc      每个链接      传递 URI, 在获取内容后可以干点啥
    rule级别：
      1. rule.beforeRuleFunc      每个链接      传递 Rule, URI, 获取内容后，执行绑定的规则方法前（如：分析字段前），用于干预内容字段分析，替换内容啥的
      2. rule.afterRuleRunc       每个链接      传递 Rule, URI, 分析字段后，执行完绑定的规则方法
      3. rule.beforeSaveFunc      每个链接      传递 Rule, URI, dataMap, 可以再这个阶段格式化内容然后再返回
      4. rule.saveFunc            每个链接      传递 Rule.ID, pk, dataMap, 返回获取到的数据字段
      5. rule.afterSaveFunc       每个链接      传递 Rule, URI, 可以在这个阶段执行通知，记录，日志等
      6. rule.fieldFilterFuncs    每个字段      传递 Field, 后置, 规则级别的全局过滤
    字段级别：
      1. field.fieldFilterFuncs   每个字段      传递 Field, 优先, 字段级别的过滤，优先于全局过滤器执行

  流程图：
    task->init->PrepareFunc->URLinitFunc->url
        url->rule->BeforeFetchFunc-fetch->AfterFetchFunc->AntiSpiderFunc->beforeRuleFunc->fetchXXX->afterRuleRunc
            fetchURL->
            fetchField->field.fieldFilterFuncs->rule.fieldFilterFuncs->save
            save->beforeSaveFunc->saveFunc->afterSaveFunc
