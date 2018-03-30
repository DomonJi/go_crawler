package pipeline

import (
	"io/ioutil"
	"strings"
		"../common/com_interfaces"
		"../common/page_items"
    "net/http"
)

type PipelineElasticsearch struct {
		url string
}

func NewPipelineElasticsearch(url string) *PipelineElasticsearch {
	  return &PipelineElasticsearch{url: url}
}

func (this * PipelineElasticsearch) Process(items *page_items.PageItems, t *com_interfaces.Task) {
    println("----------------------------------------------------------------------------------------------")
		println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
		resp, err := http.Post(this.url, "application/x-www-form-urlencoded", strings.NewReader("a=b"))
		if err != nil {
			println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			println(err)
		}
		println(string(body))
}