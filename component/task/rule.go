package task

import (
	"regexp"
	"strconv"

	"github.com/safeie/spider/common/log"
	"github.com/safeie/spider/common/util"
	"github.com/safeie/spider/component/url"
)

const (
	workFlowDrop      = iota // 丢弃
	workFlowFetchURLs        // 获取URL
	workFlowFetchRow         // 获取数据字段
	workFlowSave             // 保存
)

// Rule 任务的一个规则，不同的规则对应不同的处理流程
type Rule struct {
	task             *Task                 // 规则对应的任务
	name             string                // 规则，名称，用于区分存储的数据表
	rule             string                // 规则的字面
	re               *regexp.Regexp        // 规则的正则
	workflow         []int                 // 工作流，每一个数字，代表着一个执行方法
	pageType         int                   // 页面类型，默认 HTML网页
	forceUpdate      bool                  // 遇到采集过的页面，是否强制更新
	row              []*url.Field          // 一条数据，由多个字段组成
	pk               string                // 一条数据的主键，用于重复判断，默认为URL
	expand           bool                  // 展开单个复数字段，即：一个Row只有一个字段且该字段为数组时，展开该字段为多条数据
	fieldFilterFuncs []url.FieldFilterFunc // 字段，过滤函数，全局过滤器在字段本身过滤器之后执行
	beforeRuleFunc   BeforeRuleFunc        // 规则，前置钩子函数，在匹配到URL后，处理前
	afterRuleRunc    AfterRuleFunc         // 规则，后置钩子函数，在处理完绑定的方法后
	saveFunc         SaveFunc              // 存储，存储函数
	beforeSaveFunc   BeforeSaveFunc        // 存储，前置钩子函数
	afterSaveFunc    AfterSaveFunc         // 存储，后置钩子函数
}

// BeforeRuleFunc 规则处理前置方法
type BeforeRuleFunc func(t *Task, uri *url.URI)

// AfterRuleFunc 规则处理后置方法
type AfterRuleFunc func(t *Task, uri *url.URI)

// NewRule 创建一个新的规则
func newRule(s string, t *Task) *Rule {
	re := regexp.MustCompile(s)
	r := new(Rule)
	r.task = t
	r.rule = s
	r.re = re
	return r
}

// Task 返回规则对应的任务
func (r *Rule) Task() *Task {
	return r.task
}

// Drop 配置该配置的下一步是，丢弃
func (r *Rule) Drop() *Rule {
	r.workflow = append(r.workflow, workFlowDrop)
	return r
}

// URLs 配置该配置的下一步是，抓取URL并加入到URL列表
func (r *Rule) URLs() *Rule {
	r.workflow = append(r.workflow, workFlowFetchURLs)
	return r
}

// Row 配置该配置的下一步是，抓取字段组成一条数据
func (r *Rule) Row(fs ...*url.Field) *Rule {
	r.workflow = append(r.workflow, workFlowFetchRow)
	for i := range fs {
		if fs[i].Remote != nil {
			fs[i].Remote.SetCookie(r.task.setting.fetchOption.GetCookie())
			fs[i].Remote.SetCharset(r.task.setting.fetchOption.GetCharset())
			fs[i].Remote.SetProxy(r.task.setting.fetchOption.GetProxy())
		}
		r.row = append(r.row, fs[i])
	}
	return r
}

// PK 设置主键字段，应是数据行中的一个字段名，默认以当前页面的URL作为主键
func (r *Rule) PK(name string) *Rule {
	r.pk = name
	return r
}

// Save 配置该配置的下一步是，保存
func (r *Rule) Save() *Rule {
	r.workflow = append(r.workflow, workFlowSave)
	return r
}

// Name 获取规则名称
func (r *Rule) Name() string {
	return r.name
}

// SetName 设置规则名称
func (r *Rule) SetName(name string) *Rule {
	r.name = name
	return r
}

// SetPageType 设置页面类型，默认HTML
func (r *Rule) SetPageType(v int) *Rule {
	r.pageType = v
	return r
}

// ForceUpdate 设置符合规则的页面是否强制更新内容
func (r *Rule) ForceUpdate(v bool) *Rule {
	r.forceUpdate = v
	return r
}

// SetExpand 展开单个复数字段，即：一个Row只有一个字段且该字段为数组时，展开该字段为多条数据
func (r *Rule) SetExpand(v bool) *Rule {
	r.expand = v
	return r
}

// SetRuleFunc 设置规则的钩子函数
func (r *Rule) SetRuleFunc(b BeforeRuleFunc, a AfterRuleFunc) *Rule {
	r.beforeRuleFunc = b
	r.afterRuleRunc = a
	return r
}

// SetFieldFilterFunc 设置字段的过滤函数，全局，对所有字段生效
func (r *Rule) SetFieldFilterFunc(ff ...url.FieldFilterFunc) *Rule {
	r.fieldFilterFuncs = append(r.fieldFilterFuncs, ff...)
	return r
}

// SetSaveFunc 设置存储的钩子函数
func (r *Rule) SetSaveFunc(s SaveFunc, b BeforeSaveFunc, a AfterSaveFunc) *Rule {
	r.saveFunc = s
	r.beforeSaveFunc = b
	r.afterSaveFunc = a
	return r
}

// Match 判断一个URI是否符合该规则
func (r *Rule) Match(u string) bool {
	return r.re.MatchString(u)
}

// Run 执行规则绑定的动作
func (r *Rule) Run(u *url.URI, fetcher *FetcherPool) error {
	var err error

	// 前置方法
	if r.beforeRuleFunc != nil {
		r.beforeRuleFunc(r.task, u)
	}

	for _, w := range r.workflow {
		if w == workFlowDrop {
			break // drop
		}
		if err = r.fetch(u, fetcher); err != nil {
			log.Errorf("url fetch error: %s %v\n", u.URL, err)
			return err
		}
		switch w {
		case workFlowFetchURLs:
			urls := u.FetchURLs()
			r.task.PushURL(urls...)
		case workFlowFetchRow:
			r.parseRow(u)
		case workFlowSave:
			r.saveFields(u)
		default:
			// drop
		}
	}

	// 后置方法
	if r.afterRuleRunc != nil {
		r.afterRuleRunc(r.task, u)
	}

	return nil
}

// SaveRow 保存一条数据
func (r *Rule) SaveRow(u *url.URI, val map[string]interface{}) error {
	// before
	if r.beforeSaveFunc != nil {
		val = r.beforeSaveFunc(r, u, val)
	}

	if val == nil {
		return Errorf("数据为空")
	}

	var pk string
	if v, ok := val[r.pk]; ok {
		switch v.(type) {
		case string:
			pk = v.(string)
		case float64:
			pk = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		}
	}
	if pk == "" {
		pk = u.URL
	}

	// 保存数据
	if r.saveFunc != nil {
		if err := r.saveFunc(r.task.ID(), pk, val); err != nil {
			return err
		}
	} else {
		defaultSaveFunc(r.task.ID(), pk, val)
	}

	// after
	if r.afterSaveFunc != nil {
		r.afterSaveFunc(r, u)
	}

	return nil
}

// fetch 获取一个URI的数据并记录
func (r *Rule) fetch(u *url.URI, fetcherPool *FetcherPool) error {
	if u.Fetched {
		return nil
	}
	// 先获取URL内容
	u.PageType = r.pageType
	// 执行前置方法
	if r.task.setting.beforeFetchFunc != nil {
		r.task.setting.beforeFetchFunc(u)
	}
	// 执行重复检测
	if r.task.setting.checkRepeatFunc != nil {
		if r.task.setting.checkRepeatFunc(u) {
			return nil // 已经采集过了
		}
	}
	// 执行获取
	_, err := r.task.FetchURI(u, fetcherPool)
	// 反采集检测，无论如何都要执行，因为可能抓取错误就是反采集造成的
	if r.task.setting.antiSpiderFunc != nil && r.task.setting.antiSpiderFunc(u) {
		r.task.Logf("触发反采集策略 %s", u.URL)
		return Errorf("触发反采集策略")
	}
	// 如果不是反采集错误，判断其他错误
	if err != nil {
		return err
	}
	// 记录URL
	exists := r.task.logURL(u.URL, util.MD5Bytes(u.Body))
	if r.forceUpdate == false && exists == true {
		return nil // 无需处理，采集过了
	}
	// 执行后置方法，只在抓取成功后，并且有更新时，执行
	if r.task.setting.afterFetchFunc != nil {
		r.task.setting.afterFetchFunc(u)
	}
	return nil
}

// fields 提取该URL绑定的字段数据
func (r *Rule) parseRow(u *url.URI) {
	var err error
	for i := range r.row {
		f := r.row[i].Copy()
		if err = f.Fetch(u); err != nil {
			err = Errorf("字段[%s:%s]分析失败：%v", f.Name, f.Alias, err)
			r.task.Log(err)
		}
		// 执行全局字段过滤器，PS. 字段本身的过滤器已经优先执行
		if err == nil {
			r.filterField(f)
		}
		u.AddFields(f)
	}
}

// filterField 对字段执行过滤，含子字段
func (r *Rule) filterField(f *url.Field) {
	if len(r.fieldFilterFuncs) == 0 {
		return
	}

	children := f.Children()
	if len(children) > 0 {
		for i := range children {
			for key := range children[i] {
				r.filterField(children[i][key])
			}
		}
		return
	}

	for _, fn := range r.fieldFilterFuncs {
		fn(f)
	}
}

// saveFields 执行保存字段的动作，如果数据为空，跳过保存
// 允许before方法修复要保存的数据
func (r *Rule) saveFields(u *url.URI) {
	v := u.ExportFields()
	// 如果只有一个字段，并且设置了展开，那么久展开为多条数据
	if r.expand && len(v) == 1 {
		for _, v := range v {
			if array, ok := v.([]map[string]interface{}); ok {
				for i := range array {
					if err := r.SaveRow(u, array[i]); err != nil {
						log.Errorf("Rule.saveFields error: %s %v\n", u.URL, err)
					}
				}
			}
			break
		}
		return
	}

	// 不展开，直接保存这条数据
	if v != nil {
		if err := r.SaveRow(u, v); err != nil {
			log.Errorf("Rule.saveFields error: %s %v\n", u.URL, err)
		}
	}

}
