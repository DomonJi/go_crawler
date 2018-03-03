package page_processer

import (
    "../common/page"
)

type PageProcesser interface {
    Process(p *page.Page)
    Finish()
}
