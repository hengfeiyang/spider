package main

import (
	"fmt"
	"strings"

	"github.com/safeie/spider/component/task"
	"github.com/safeie/spider/component/url"
)

func main() {
	// create a new task
	t := task.New("1", "zhihu", "https://www.zhihu.com", "")
	// set init urls
	t.SetURLinitFunc(func() []string {
		return []string{"https://www.zhihu.com/question/29252365"}
	})
	// enable webkit engine
	t.EnableJS(true)
	// set webkit delay time before fetch
	t.SetRenderDelay(1)
	// set goruntime num
	t.SetRoutineNum(1)
	// set continue
	t.SetErrorContinue(true)
	// set callback handle author url
	t.SetFetchFunc(nil, nil, func(u *url.URI) {
		content := string(u.Body)
		content = strings.Replace(content, "https://www.zhihu.com/people/", "", -1)
		u.ResetBody([]byte(content))
	})

	// prepare page fields
	// use html document selector query
	authorName := t.NewField("author_name", "author name").
		SetMatchRule(url.MatchTypeSubString, `<meta itemprop="name" content="(*)"/>`)
	// this is a remote field, it fetch data from the other page
	authorInfo := t.NewField("author_info", "author info").
		SetMatchRule(url.MatchTypeSubString, `<meta itemprop="url" content="(*)"/>`).
		// {{.}} will replace with f.Value, such as author.id, author.name
		SetRemote(url.NewRemote(t, url.PageTypeHTML, "https://www.zhihu.com/people/{{.}}")).
		// set children field in the remote page
		SetChildren(
			t.NewField("desc", "description").SetMatchRule(url.MatchTypeSubString, `,"description":"(*)"`),
		)

	items := t.NewField("items", "answers").
		SetMatchRule(url.MatchTypeSelector, "#QuestionAnswers-answers > div > div > div > div:nth-child(2) > div > div").
		SetRepeat(true).
		SetChildren(authorName, authorInfo)

	// use regexp match url rule
	t.Rule("https://www.zhihu.com/question/29252365$").
		URLs().     // collect matched url
		Row(items). // parse page fields
		SetSaveFunc(func(taskID, pk string, val map[string]interface{}) error {
			data, _ := val["items"].([]map[string]interface{})
			for i, row := range data {
				fmt.Println(i, row["author_name"], row["author_info"])
				fmt.Println("-----------------")
			}
			return nil
		}, nil, nil).Save()

	// run
	fmt.Printf("done with: %v\n", t.Run())
}
