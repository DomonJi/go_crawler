package scheduler

import (
    "container/list"
    "../common/request"
    "../bloom"
    "sync"
)

type QueueScheduler struct {
    locker *sync.Mutex
    rm     bool
    filter *bloom.BloomFilter
    queue  *list.List
}

func NewQueueScheduler(rmDuplicate bool) *QueueScheduler {
    queue := list.New()
    filter := bloom.NewWithEstimates(100000, 0.001)
    locker := new(sync.Mutex)
    return &QueueScheduler{rm: rmDuplicate, queue: queue, locker: locker, filter: filter}
}

func (this *QueueScheduler) Push(requ *request.Request) {
    this.locker.Lock()
    url := requ.GetUrl()
    if this.rm {
        if this.filter.TestString(url) {
            this.locker.Unlock()
            return
        }
    }
    this.queue.PushBack(requ)
    if this.rm {
        this.filter.AddString(url)
    }
    this.locker.Unlock()
}

func (this *QueueScheduler) Poll() *request.Request {
    this.locker.Lock()
    if this.queue.Len() <= 0 {
        this.locker.Unlock()
        return nil
    }
    e := this.queue.Front()
    requ := e.Value.(*request.Request)
    this.queue.Remove(e)
    this.locker.Unlock()
    return requ
}

func (this *QueueScheduler) Count() int {
    this.locker.Lock()
    len := this.queue.Len()
    this.locker.Unlock()
    return len
}
