package page_items

import (
    "../request"
)

type PageItems struct {

    req *request.Request

    items map[string]string

    skip bool
}

func NewPageItems(req *request.Request) *PageItems {
    items := make(map[string]string)
    return &PageItems{req: req, items: items, skip: false}
}

func (this *PageItems) GetRequest() *request.Request {
    return this.req
}

func (this *PageItems) AddItem(key string, item string) {
    this.items[key] = item
}

func (this *PageItems) GetItem(key string) (string, bool) {
    t, ok := this.items[key]
    return t, ok
}

func (this *PageItems) GetAll() map[string]string {
    return this.items
}

func (this *PageItems) SetSkip(skip bool) *PageItems {
    this.skip = skip
    return this
}

func (this *PageItems) GetSkip() bool {
    return this.skip
}
