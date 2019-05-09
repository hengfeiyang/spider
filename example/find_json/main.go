package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/safeie/spider/component/task"
	"github.com/safeie/spider/component/url"
)

func main() {
	// create a new task
	t := task.New("1", "dongde", "https://www.idongde.com", "")
	// set crawl sleep time
	t.SetInterval(1000)
	// set goruntime num
	t.SetRoutineNum(1)
	// set init urls
	t.SetURLinitFunc(func() []string {
		urls := make([]string, 0)
		for i := 1; i < 10; i++ {
			urls = append(urls, fmt.Sprintf("https://www.idongde.com/index/page?page=%d&size=18", i))
		}
		return urls
	})

	// i can the queue before quit
	t.SetBeforeQuitFunc(func(taskID string, queue []string) {
		fmt.Printf("i received the remain queue from task: %s, queue is:\n", taskID)
		for i, item := range queue {
			fmt.Printf("%5d\t%s\n", i, item)
		}
	})

	// set json field rule
	// use simply json parser
	data := t.NewField("data", "data").SetMatchRule(url.MatchTypeJSONPath, "$.data.data").
		SetRepeat(true). // data is array
		SetChildren(     // row is object
			t.NewField("id", "id").SetMatchRule(url.MatchTypeJSONPath, "$.id"),
			t.NewField("title", "page title").SetMatchRule(url.MatchTypeJSONPath, "$.title"),
		)

	var num uint32
	t.Rule("https://www.idongde.com/index/page*").
		SetName("video rule"). // set rule name
		URLs().                // collect matched url
		Row(data).             // parse page fields
		SetSaveFunc(func(taskID, pk string, val map[string]interface{}) error {
			atomic.AddUint32(&num, 1)
			fmt.Printf("pk: %s\n-------------------------\n", pk)
			data := val["data"]
			dataSlice, _ := data.([]map[string]interface{})
			for j, row := range dataSlice {
				fmt.Printf("%5d\tid: %6.0f\ttitle: %s\n", j, row["id"], row["title"])
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
