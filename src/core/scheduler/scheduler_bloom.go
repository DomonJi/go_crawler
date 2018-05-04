package scheduler

import (
    "../common/request"
    "../bloom"
)

type BloomScheduler struct {
    queue chan *request.Request
    filter *bloom.BloomFilter
}

func NewBloomScheduler() *BloomScheduler {
    ch := make(chan *request.Request, 4096)
    filter := bloom.NewWithEstimates(100000, 0.001)
    return &BloomScheduler{queue: ch, filter: filter}
}

func (this *BloomScheduler) Push(requ *request.Request) {
    url := requ.GetUrl()
    if this.filter.TestString(url) {
        return
    }
    this.queue <- requ
    this.filter.AddString(url)
}

func (this *BloomScheduler) Poll() *request.Request {
    if len(this.queue) == 0 {
        return nil
    } else {
        return <-this.queue
    }
}

func (this *BloomScheduler) Count() int {
    return len(this.queue)
}
