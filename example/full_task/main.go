package main

import (
	"fmt"

	"github.com/safeie/spider/component/task"
	"github.com/safeie/spider/component/url"
)

func main() {
	spider()
	select {}
}

func spider() {
	t := task.New("1", "golang blog", "https://blog.golang.org")
	t.SetURLinitFunc(func() []string {
		return []string{"https://blog.golang.org/index"}
	})
	t.SetFetchFunc(nil, func(uri *url.URI) {
		fmt.Println("before", uri)
	}, func(uri *url.URI) {
		fmt.Println("after", uri)
	})

	title := t.NewField("title", "标题").SetMatchRule(url.MatchTypeSelector, "#content > div > h3 > a")
	content := t.NewField("content", "内容").SetMatchRule(url.MatchTypeSelector, "#content > div").
		SetFilterFunc(func(f *url.Field) {
			f.Remove(`<h3 class="title">(*)</h3>`)
			f.Remove(`<p class="date">(*)</p>`)
		}).
		SetFixURL(true)

	t.Rule("https://blog.golang.org/*").
		URLs().
		Row(title, content).
		SetSaveFunc(func(taskID, pk string, val map[string]interface{}) error {
			fmt.Printf("taskID: %s\npk: %s\ntitle: %s\ncontent:\n%s\n\n", taskID, pk, val["title"], val["content"])
			return nil
		}, nil, nil).Save()

	err := t.Run()
	fmt.Printf("run: %v\n", err)
}
