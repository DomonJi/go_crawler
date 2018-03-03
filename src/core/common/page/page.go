package page

import (
    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "../mlog"
    "../page_items"
    "../request"
    "net/http"
    "strings"
    //"fmt"
)

type Page struct {
    isfail   bool
    errormsg string

    req *request.Request

    body string

    header  http.Header
    cookies []*http.Cookie

    docParser *goquery.Document

    jsonMap *simplejson.Json

    pItems *page_items.PageItems

    targetRequests []*request.Request
}

func NewPage(req *request.Request) *Page {
    return &Page{pItems: page_items.NewPageItems(req), req: req}
}

func (this *Page) SetHeader(header http.Header) {
    this.header = header
}

func (this *Page) GetHeader() http.Header {
    return this.header
}

func (this *Page) SetCookies(cookies []*http.Cookie) {
    this.cookies = cookies
}

func (this *Page) GetCookies() []*http.Cookie {
    return this.cookies
}

func (this *Page) IsSucc() bool {
    return !this.isfail
}

func (this *Page) Errormsg() string {
    return this.errormsg
}

func (this *Page) SetStatus(isfail bool, errormsg string) {
    this.isfail = isfail
    this.errormsg = errormsg
}

func (this *Page) AddField(key string, value string) {
    this.pItems.AddItem(key, value)
}

func (this *Page) GetPageItems() *page_items.PageItems {
    return this.pItems
}

func (this *Page) SetSkip(skip bool) {
    this.pItems.SetSkip(skip)
}

func (this *Page) GetSkip() bool {
    return this.pItems.GetSkip()
}

func (this *Page) SetRequest(r *request.Request) *Page {
    this.req = r
    return this
}

func (this *Page) GetRequest() *request.Request {
    return this.req
}

func (this *Page) GetUrlTag() string {
    return this.req.GetUrlTag()
}

func (this *Page) AddTargetRequest(url string, respType string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequest(url, respType, "", "GET", "", nil, nil, nil, nil))
    return this
}

func (this *Page) AddTargetRequests(urls []string, respType string) *Page {
    for _, url := range urls {
        this.AddTargetRequest(url, respType)
    }
    return this
}

func (this *Page) AddTargetRequestWithProxy(url string, respType string, proxyHost string) *Page {

    this.targetRequests = append(this.targetRequests, request.NewRequestWithProxy(url, respType, "", "GET", "", nil, nil, proxyHost, nil, nil))
    return this
}

func (this *Page) AddTargetRequestsWithProxy(urls []string, respType string, proxyHost string) *Page {
    for _, url := range urls {
        this.AddTargetRequestWithProxy(url, respType, proxyHost)
    }
    return this
}

func (this *Page) AddTargetRequestWithHeaderFile(url string, respType string, headerFile string) *Page {
    this.targetRequests = append(this.targetRequests, request.NewRequestWithHeaderFile(url, respType, headerFile))
    return this
}

func (this *Page) AddTargetRequestWithParams(req *request.Request) *Page {
    this.targetRequests = append(this.targetRequests, req)
    return this
}

func (this *Page) AddTargetRequestsWithParams(reqs []*request.Request) *Page {
    for _, req := range reqs {
        this.AddTargetRequestWithParams(req)
    }
    return this
}

func (this *Page) GetTargetRequests() []*request.Request {
    return this.targetRequests
}

func (this *Page) SetBodyStr(body string) *Page {
    this.body = body
    return this
}

func (this *Page) GetBodyStr() string {
    return this.body
}

func (this *Page) SetHtmlParser(doc *goquery.Document) *Page {
    this.docParser = doc
    return this
}

func (this *Page) GetHtmlParser() *goquery.Document {
    return this.docParser
}

func (this *Page) ResetHtmlParser() *goquery.Document {
    r := strings.NewReader(this.body)
    var err error
    this.docParser, err = goquery.NewDocumentFromReader(r)
    if err != nil {
        mlog.LogInst().LogError(err.Error())
        panic(err.Error())
    }
    return this.docParser
}

func (this *Page) SetJson(js *simplejson.Json) *Page {
    this.jsonMap = js
    return this
}

func (this *Page) GetJson() *simplejson.Json {
    return this.jsonMap
}
