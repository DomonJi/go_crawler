package pipeline

import (
	  "context"
		"../common/com_interfaces"
		"../common/page_items"
		"github.com/olivere/elastic"
)

type PipelineElasticsearch struct {
	client *elastic.Client
}

func NewPipelineElasticsearch(client *elastic.Client) *PipelineElasticsearch {
	  return &PipelineElasticsearch{client: client}
}

func (this *PipelineElasticsearch) Process(items *page_items.PageItems, t com_interfaces.Task) {
    println("----------------------------------------------------------------------------------------------")
		println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
		_, err := this.client.Index().Index("spider").Type("doc").BodyJson(items.GetAll()).Do(context.Background())
		if err != nil {
			println(err)
		}
}