package main

import (
	"fmt"
	"sync/atomic"

	"github.com/safeie/spider/component/task"
	"github.com/safeie/spider/component/url"
)

func main() {
	// create a new task
	t := task.New("1", "golang blog", "https://blog.golang.org", "")
	// set crawl sleep time
	t.SetInterval(100)
	// set rand useragent
	t.SetUserAgentPool("ua1", "ua2")
	// set goruntime num
	t.SetRoutineNum(2)
	// set continue
	t.SetErrorContinue(true)
	// set global cookie
	t.SetPrepareFunc(func(t *task.Task) {
		_, cookie, err := t.FetchURL("https://blog.golang.org/index")
		if err == nil {
			t.SetCookie(cookie)
		}
	})
	// set init urls
	t.SetURLinitFunc(func() []string {
		return []string{"https://blog.golang.org/index"}
	})
	// collect urls
	t.Rule("https://blog.golang.org/index").URLs()

	// prepare page fields
	// use html document selector query
	title := t.NewField("title", "title").SetMatchRule(url.MatchTypeSelector, "#content > div > h3 > a")
	content := t.NewField("content", "content").SetMatchRule(url.MatchTypeSelector, "#content > div").
		SetFilterFunc(func(f *url.Field) {
			f.Remove(`<h3 class="title">(*)</h3>`)
			f.Remove(`<p class="date">(*)</p>`)
		}).
		SetFixURL(true)

	var num uint32
	// use regexp match url rule
	t.Rule("https://blog.golang.org/*").
		SetName("blog paper"). // set rule name
		URLs().                // collect matched url
		Row(title, content).   // parse page fields
		SetSaveFunc(func(taskID, pk string, val map[string]interface{}) error {
			atomic.AddUint32(&num, 1)
			fmt.Printf("%5d taskID: %s, \n->    pk: %s\n-> title: %s\n", num, taskID, pk, val["title"])
			return nil
		}, nil, nil).
		Save() // save rule

	// run
	fmt.Printf("done with: %v\n", t.Run())
}
