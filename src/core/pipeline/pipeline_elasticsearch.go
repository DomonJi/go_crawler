package pipeline

import (
		"../common/com_interfaces"
		"../common/page_items"
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
}