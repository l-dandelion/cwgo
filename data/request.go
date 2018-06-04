package data

import (
	"net/http"
)

type Request struct {
	httpReq *http.Request
	depth   int                    // 抓取深度
	proxy   string                 // 代理信息
	Extra   map[string]interface{} // 额外信息
}

/*
 * New an instance of Request
 */
func NewRequest(httpReq *http.Request, extras ...map[string]interface{}) *Request {
	var extra map[string]interface{}
	if len(extras) != 0 {
		extra = extras[0]
	} else {
		extra = map[string]interface{}{}
	}
	return &Request{
		httpReq: httpReq,
		Extra:   extra,
	}
}

func (req *Request) Valid() bool {
	return req != nil && req.httpReq != nil
}

/*Get*/

func (req *Request) HTTPReq() *http.Request {
	return req.httpReq
}

func (req *Request) Depth() int {
	return req.depth
}

func (req *Request) Proxy() string {
	return req.proxy
}

func (req *Request) SetExtra(key string, val interface{}) {
	req.Extra[key] = val
}

/*Set*/

func (req *Request) SetDepth(depth int) {
	req.depth = depth
}

func (req *Request) SetHeader(key, val string) {
	req.httpReq.Header.Set(key, val)
}

//cookie
func (req *Request) AddCookie(key, value string) {
	c := &http.Cookie{
		Name:  key,
		Value: value,
	}
	req.httpReq.AddCookie(c)
}

//ua
func (req *Request) SetUserAgent(ua string) {
	req.SetHeader("User-Agent", ua)
}

//referer
func (req *Request) SetReferer(referer string) {
	req.SetHeader("referer", referer)
}

//proxy
func (req *Request) SetProxy(proxy string) {
	req.proxy = proxy
}


