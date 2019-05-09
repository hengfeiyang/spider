// Package wendu 文都网执业药师采集
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/safeie/spider/common/log"
	"github.com/safeie/spider/component/parser"
	"github.com/safeie/spider/component/url"
)

const storeFile = "store.csv"

type docItem struct {
	Title       string
	DownloadURL string
}

func main() {
	crawl()
}

func crawl() {
	listPrefix := "http://www.wendu.com/index.php?m=content&c=index&a=lists&siteid=1&catid=303&page="
	listNum := 38

	items := make([]*docItem, 0)
	for i := 1; i <= listNum; i++ {
		time.Sleep(time.Second)
		pageURL := fmt.Sprintf("%s%d", listPrefix, i)
		fmt.Println(pageURL)

		docs, err := crawlPage(pageURL)
		if err != nil {
			log.Fatalln(err)
		}
		items = append(items, docs...)
	}

	// 写入文件
	file, err := os.OpenFile(storeFile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 777)
	if err != nil {
		log.Fatalln(err)
	}
	for i, item := range items {
		file.WriteString(fmt.Sprintf("%d,%s,%s\n", i, item.Title, item.DownloadURL))
	}
	file.Close()
	fmt.Println("crawl ok, total docs ", len(items))
}

// crawlPage 分析一个列表页，返回列表中每个页面的标题和文档下载地址
func crawlPage(uri string) ([]*docItem, error) {
	remote := url.NewRemote(log.DefaultLogger.(log.SimpleLogger), url.PageTypeHTML, uri)
	pageRet, err := remote.FetchURI("")
	if err != nil {
		return nil, err
	}
	pageParser, err := parser.NewSubstring(pageRet.Body)
	if err != nil {
		return nil, err
	}
	docUrls := pageParser.MatchAll(`date-load-btn" href="(*)">`)

	items := make([]*docItem, 0)
	for _, docURL := range docUrls {
		docURL = pageRet.FixURL(docURL)
		time.Sleep(time.Millisecond * 10)
		fmt.Println(docURL)

		docRemote := url.NewRemote(log.DefaultLogger.(log.SimpleLogger), url.PageTypeHTML, docURL)
		docRet, err := docRemote.FetchURI("")
		if err != nil {
			return nil, err
		}
		docParser, err := parser.NewSubstring(docRet.Body)
		if err != nil {
			return nil, err
		}
		// 分析字段
		doc := new(docItem)
		doc.Title = docParser.Match(`<p class="subtitle">(*)</p>`)
		doc.DownloadURL = docParser.Match(`article-download" href="(*)">`)
		items = append(items, doc)
	}

	return items, nil
}
