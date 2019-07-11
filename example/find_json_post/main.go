package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/safeie/spider/component/task"
	"github.com/safeie/spider/component/url"
)

func main() {
	// create a new task
	t := task.New("1", "bevol", "https://internal.bevol.cn/", "")
	// set crawl sleep time
	t.SetInterval(1000)
	// set goruntime num
	t.SetRoutineNum(2)
	t.SetMethod(task.MethodPOST)
	// set init urls
	t.SetURLinitFunc(func() []string {
		urls := make([]string, 0)
		for i := 1; i <= 3; i++ {
			urls = append(urls, fmt.Sprintf("https://internal.bevol.cn/good/article/findByList?pager=%d", i))
		}
		return urls
	})
	t.SetParam("pageSize", "10")
	t.SetParam("hidden", "0")
	t.SetParam("pager", "1")
	t.SetFetchFunc(nil, func(uri *url.URI) {
		pos := strings.Index(uri.URL, "=")
		pager := uri.URL[pos+1:]
		uri.SetParam("pager", pager)
	}, nil)

	// i can the queue before quit
	t.SetBeforeQuitFunc(func(taskID string, queue []string) {
		fmt.Printf("i received the remain queue from task: %s, queue is:\n", taskID)
		for i, item := range queue {
			fmt.Printf("%5d\t%s\n", i, item)
		}
	})

	// set json field rule
	// use simply json parser
	title := t.NewField("title", "page title").SetMatchRule(url.MatchTypeJSONPath, "$.title")
	content := t.NewField("content", "page content").SetMatchRule(url.MatchTypeJSONPath, "$.mid").
		SetRemote(url.NewRemote(t, url.PageTypeJSON, "https://internal.bevol.cn/good/article/article").
			SetMethod(task.MethodPOST).
			SetParam("mid", "{{.}}")).
		SetChildren(
			t.NewField("mid", "mid").SetMatchRule(url.MatchTypeJSONPath, "$.result.entity.mid"),
			t.NewField("data", "data").SetMatchRule(url.MatchTypeJSONPath, "$.result.entityInfo"),
		)
	data := t.NewField("data", "data").SetMatchRule(url.MatchTypeJSONPath, "$.result.list").
		SetRepeat(true). // data is array
		SetChildren(     // row is object
			title,
			content,
		)

	t.Rule("https://internal.bevol.cn/good/article/findByList*").
		SetName("content rule"). // set rule name
		Row(data).               // parse page fields
		SetSaveFunc(func(taskID, pk string, val map[string]interface{}) error {
			fmt.Printf("pk: %s\n-------------------------\n", pk)
			data := val["data"]
			dataSlice, _ := data.([]map[string]interface{})
			for j, row := range dataSlice {
				content := row["content"].(map[string]interface{})
				fmt.Printf("%5d\tmid: %v\ttitle: %s\n", j+1, content["mid"], row["title"])
			}
			return nil
		}, nil, nil).Save()

	// trap signal
	close := make(chan struct{})
	trapSignal(close)
	go func(close chan struct{}) {
		<-close
		t.Stop()
	}(close)

	// run
	fmt.Printf("done with: %v\n", t.Run())
}

func trapSignal(close chan struct{}) {
	sch := make(chan os.Signal, 10)
	signal.Notify(sch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT,
		syscall.SIGHUP, syscall.SIGSTOP, syscall.SIGQUIT)
	go func(ch <-chan os.Signal) {
		sig := <-ch
		log.Println("signal recieved " + sig.String() + ", at: " + time.Now().String())
		close <- struct{}{}
		if sig == syscall.SIGHUP {
			log.Println("restart now...")
			procAttr := new(os.ProcAttr)
			procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
			procAttr.Dir = os.Getenv("PWD")
			procAttr.Env = os.Environ()
			process, err := os.StartProcess(os.Args[0], os.Args, procAttr)
			if err != nil {
				log.Println("restart process failed:" + err.Error())
				return
			}
			waitMsg, err := process.Wait()
			if err != nil {
				log.Println("restart wait error:" + err.Error())
			}
			log.Println(waitMsg)
		}
	}(sch)
}
