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

```
task->init->PrepareFunc->URLinitFunc->url
  url->rule->BeforeFetchFunc->CheckRepeatFunc->fetch->AfterFetchFunc->AntiSpiderFunc->beforeRuleFunc->parse->afterRuleFunc
    fetchURL->
    fetchField->field.fieldFilterFuncs->rule.fieldFilterFuncs->save
    save->beforeSaveFunc->saveFunc->afterSaveFunc
```

### callback

* task level:
  1. task.PrepareFunc         once         pass Task, can do something initialize, eg: setCookie
  2. task.URLinitFunc         once         pass None, receive initialize urls, it should multiple
  3. task.AntiSpiderFunc      per uri      pass Uri,  check crawl behaviour is trigger the anti spider rule
  4. task.BeforeFetchFunc     per uri      pass Uri,  you can do something before fetch url content
  5. task.CheckRepeatFunc     per uri      pass Uri,  you can check is repeated, if return true will skip the url
  6. task.AfterFetchFunc      per uri      pass Uri,  you can do something after fetch url content
  7. task.BeforeQuitFunc      once         pass taskid and url queue, you can storage saved queue before quit, it can be used for next init urls
* rule level:
  1. rule.beforeRuleFunc      per uri      pass Rule, URI, after fetch and before excute rule functions, you can change the page content or something else
  2. rule.afterRuleRunc       per uri      pass Rule, URI, after execute rule parse functions
  3. rule.beforeSaveFunc      per uri      pass Rule, URI, dataMap, before save data, you have a change to filter the data
  4. rule.saveFunc            per uri      pass Rule.ID, pk, dataMap, returns the data to you, you should save it
  5. rule.afterSaveFunc       per uri      pass Rule, URI, after save, you can log or dispath the success message
  6. rule.fieldFilterFuncs    per field    pass Field, a rule level global field filter, execute on each field
* filed level:
  1. field.fieldFilterFuncs   per filed    pass Field, filter the filed value, this execute before the rule level filter 
