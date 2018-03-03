package pipeline

import (
    "../common/com_interfaces"
    "../common/page_items"
)

type Pipeline interface {
    Process(items *page_items.PageItems, t com_interfaces.Task)
}

type CollectPipeline interface {
    Pipeline

    GetCollected() []*page_items.PageItems
}
