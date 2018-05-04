package main

import (
		"fmt"
		"context"
		"regexp"
    "../core/common/page"
    "../core/pipeline"
    "../core/spider"
    "../core/scheduler"
    "strings"
    "github.com/PuerkitoBio/goquery"
	  "github.com/olivere/elastic"
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
				match, _ := regexp.MatchString("/item/", href)
				if match {
					urls = append(urls, "https://baike.baidu.com" + href)
				}
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
		p.AddField("url", p.GetRequest().GetUrl())
}

func (this *MyPageProcesser) Finish() {
    fmt.Printf("TODO:before end spider \r\n")
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}
	exist, err := client.IndexExists("spider").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"name":{
					"type":"text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"summary":{
					"type":"text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				},
				"url":{
					"type":"text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_max_word"
				}
			}
		}
	}
}
`
		_, err := client.CreateIndex("spider").Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	spider.NewSpider(NewMyPageProcesser(), "baidu_baike_spider").
		SetScheduler(scheduler.NewQueueScheduler(true)).
		AddUrl("https://baike.baidu.com/view/1628025.htm?fromtitle=http&fromid=243074&type=syn", "html").
		AddPipeline(pipeline.NewPipelineElasticsearch(client)).
		SetSleepTime("rand", 500, 1000).
		SetThreadnum(8).
		Run()
}