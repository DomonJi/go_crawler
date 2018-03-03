package downloader

import (
    "bytes"

    "github.com/PuerkitoBio/goquery"
    "github.com/bitly/go-simplejson"
    "../common/mlog"
    "../common/page"
    "../common/request"
    "../common/util"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "golang.org/x/net/html/charset"
    "strings"
	"compress/gzip"
)

type HttpDownloader struct {
}

func NewHttpDownloader() *HttpDownloader {
    return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *request.Request) *page.Page {
    var mtype string
    var p = page.NewPage(req)
    mtype = req.GetResponceType()
    switch mtype {
    case "html":
        return this.downloadHtml(p, req)
    case "json":
        fallthrough
    case "jsonp":
        return this.downloadJson(p, req)
    case "text":
        return this.downloadText(p, req)
    default:
        mlog.LogInst().LogError("error request type:" + mtype)
    }
    return p
}

func (this *HttpDownloader) changeCharsetEncodingAuto(contentTypeStr string, sor io.ReadCloser) string {
    var err error
    destReader, err := charset.NewReader(sor, contentTypeStr)

    if err != nil {
        mlog.LogInst().LogError(err.Error())
        destReader = sor
    }

    var sorbody []byte
    if sorbody, err = ioutil.ReadAll(destReader); err != nil {
        mlog.LogInst().LogError(err.Error())
    }
    bodystr := string(sorbody)

    return bodystr
}

func (this *HttpDownloader) changeCharsetEncodingAutoGzipSupport(contentTypeStr string, sor io.ReadCloser) string {
	var err error
	gzipReader, err := gzip.NewReader(sor)
	if err != nil {
		mlog.LogInst().LogError(err.Error())
		return ""
	}
	defer gzipReader.Close()
	destReader, err := charset.NewReader(gzipReader, contentTypeStr)

	if err != nil {
		mlog.LogInst().LogError(err.Error())
		destReader = sor
	}

	var sorbody []byte
	if sorbody, err = ioutil.ReadAll(destReader); err != nil {
		mlog.LogInst().LogError(err.Error())
	}
	bodystr := string(sorbody)

	return bodystr
}

func connectByHttp(p *page.Page, req *request.Request) (*http.Response, error) {
    client := &http.Client{
        CheckRedirect: req.GetRedirectFunc(),
    }

    httpreq, err := http.NewRequest(req.GetMethod(), req.GetUrl(), strings.NewReader(req.GetPostdata()))
    if header := req.GetHeader(); header != nil {
        httpreq.Header = req.GetHeader()
    }

    if cookies := req.GetCookies(); cookies != nil {
        for i := range cookies {
            httpreq.AddCookie(cookies[i])
        }
    }

    var resp *http.Response
    if resp, err = client.Do(httpreq); err != nil {
        if e, ok := err.(*url.Error); ok && e.Err != nil && e.Err.Error() == "normal" {
        } else {
            mlog.LogInst().LogError(err.Error())
            p.SetStatus(true, err.Error())
            return nil, err
        }
    }

    return resp, nil
}

func connectByHttpProxy(p *page.Page, in_req *request.Request) (*http.Response, error) {
    request, _ := http.NewRequest("GET", in_req.GetUrl(), nil)
    proxy, err := url.Parse(in_req.GetProxyHost())
    if err != nil {
        return nil, err
    }
    client := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxy),
        },
    }
    resp, err := client.Do(request)
    if err != nil {
        return nil, err
    }
    return resp, nil

}

func (this *HttpDownloader) downloadFile(p *page.Page, req *request.Request) (*page.Page, string) {
    var err error
    var urlstr string
    if urlstr = req.GetUrl(); len(urlstr) == 0 {
        mlog.LogInst().LogError("url is empty")
        p.SetStatus(true, "url is empty")
        return p, ""
    }

    var resp *http.Response

    if proxystr := req.GetProxyHost(); len(proxystr) != 0 {
        //fmt.Print("HttpProxy Enter ",proxystr,"\n")
        resp, err = connectByHttpProxy(p, req)
    } else {
        //fmt.Print("Http Normal Enter \n",proxystr,"\n")
        resp, err = connectByHttp(p, req)
    }

    if err != nil {
        return p, ""
    }

    //b, _ := ioutil.ReadAll(resp.Body)
    //fmt.Printf("Resp body %v \r\n", string(b))

    p.SetHeader(resp.Header)
    p.SetCookies(resp.Cookies())

	var bodyStr string
	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyStr = this.changeCharsetEncodingAutoGzipSupport(resp.Header.Get("Content-Type"), resp.Body)
	} else {
		bodyStr = this.changeCharsetEncodingAuto(resp.Header.Get("Content-Type"), resp.Body)
	}
    defer resp.Body.Close()
    return p, bodyStr
}

func (this *HttpDownloader) downloadHtml(p *page.Page, req *request.Request) *page.Page {
    var err error
    p, destbody := this.downloadFile(p, req)
    //fmt.Printf("Destbody %v \r\n", destbody)
    if !p.IsSucc() {
        //fmt.Print("Page error \r\n")
        return p
    }
    bodyReader := bytes.NewReader([]byte(destbody))

    var doc *goquery.Document
    if doc, err = goquery.NewDocumentFromReader(bodyReader); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    var body string
    if body, err = doc.Html(); err != nil {
        mlog.LogInst().LogError(err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    p.SetBodyStr(body).SetHtmlParser(doc).SetStatus(false, "")

    return p
}

func (this *HttpDownloader) downloadJson(p *page.Page, req *request.Request) *page.Page {
    var err error
    p, destbody := this.downloadFile(p, req)
    if !p.IsSucc() {
        return p
    }

    var body []byte
    body = []byte(destbody)
    mtype := req.GetResponceType()
    if mtype == "jsonp" {
        tmpstr := util.JsonpToJson(destbody)
        body = []byte(tmpstr)
    }

    var r *simplejson.Json
    if r, err = simplejson.NewJson(body); err != nil {
        mlog.LogInst().LogError(string(body) + "\t" + err.Error())
        p.SetStatus(true, err.Error())
        return p
    }

    p.SetBodyStr(string(body)).SetJson(r).SetStatus(false, "")

    return p
}

func (this *HttpDownloader) downloadText(p *page.Page, req *request.Request) *page.Page {
    p, destbody := this.downloadFile(p, req)
    if !p.IsSucc() {
        return p
    }

    p.SetBodyStr(destbody).SetStatus(false, "")
    return p
}
