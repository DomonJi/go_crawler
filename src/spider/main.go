package main

import (
    "fmt"
    "../core/common/page"
    "../core/pipeline"
    "../core/spider"
    "../core/scheduler"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
    return &MyPageProcesser{}
}

func (this *MyPageProcesser) Process(p *page.Page) {
    if !p.IsSucc() {
        println(p.Errormsg())
        return
    }

		query := p.GetHtmlParser()
		var urls []string
		query.Find("div.para a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				urls = append(urls, "https://baike.baidu.com" + href)
		})
		p.AddTargetRequests(urls, "html")

    name := query.Find(".lemmaWgt-lemmaTitle-title h1").Text()
    name = strings.Trim(name, " \t\n")

    summary := query.Find(".lemma-summary .para").Text()
		summary = strings.Trim(summary, " \t\n")
		
		para := ""

		query.Find("div.main-content .para").Each(func(i int, s *goquery.Selection) {
			para = para + strings.Trim(s.Text(), " \t\n")
		})

    p.AddField("name", name)
		p.AddField("summary", summary)
}

func (this *MyPageProcesser) Finish() {
    fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	spider.NewSpider(NewMyPageProcesser(), "baidu_baike_spider").
		SetScheduler(scheduler.NewQueueScheduler(true)).
		AddUrl("https://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html").
		AddPipeline(pipeline.NewPipelineConsole()).
		SetSleepTime("rand", 500, 1000).
		SetThreadnum(8).
		Run()
}