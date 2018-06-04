package data

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cwgo/lib/reader"
	"golang.org/x/net/html/charset"
)

type Response struct {
	req      *Request          //响应对应的请求
	httpResp *http.Response    //请求对应的http响应
	text     []byte            //body的[]byte类型
	dom      *goquery.Document //dom结构
}

/*
 * New an instance of Response
 */
func NewResponse(req *Request, httpResp *http.Response) *Response {
	return &Response{req: req, httpResp: httpResp}
}

func (resp *Response) Valid() bool {
	return resp != nil && resp.httpResp != nil
}

/*Get*/

func (resp *Response) Request() *Request {
	return resp.req
}

func (resp *Response) HTTPRequest() *http.Request {
	return resp.req.HTTPReq()
}

func (resp *Response) HTTPResp() *http.Response {
	return resp.httpResp
}

func (resp *Response) Depth() int {
	return resp.req.Depth()
}

func (resp *Response) GetText() ([]byte, error) {
	if resp.text != nil {
		return resp.text, nil
	}
	multiReader, err := reader.NewMultipleReader(resp.httpResp.Body)
	resp.httpResp.Body = multiReader.Reader()
	defer func() {
		resp.httpResp.Body.Close()
		resp.httpResp.Body = multiReader.Reader()
	}()
	var contentType, pageEncode string

	// read firstly content-type from response header
	contentType = resp.httpResp.Header.Get("Content-Type")
	if _, params, err := mime.ParseMediaType(contentType); err == nil {
		if cs, ok := params["charset"]; ok {
			pageEncode = strings.ToLower(strings.TrimSpace(cs))
		}
	}

	// read content-type from request header
	if len(pageEncode) == 0 {
		contentType = resp.httpResp.Request.Header.Get("Content-Type")
		if _, params, err := mime.ParseMediaType(contentType); err == nil {
			if cs, ok := params["charset"]; ok {
				pageEncode = strings.ToLower(strings.TrimSpace(cs))
			}
		}
	}

	switch pageEncode {
	case "utf8", "utf-8", "unicode-1-1-utf-8":
	default:
		// get converter to utf-8
		// Charset auto determine. Use golang.org/x/net/html/charset. Get response body and change it to utf-8
		var destReader io.Reader

		if len(pageEncode) == 0 {
			destReader, err = charset.NewReader(resp.httpResp.Body, "")
		} else {
			destReader, err = charset.NewReaderLabel(pageEncode, resp.httpResp.Body)
		}

		if err == nil {
			resp.text, err = ioutil.ReadAll(destReader)
			if err == nil {
				return resp.text, err
			}
		}

	}
	resp.text, err = ioutil.ReadAll(resp.httpResp.Body)

	return resp.text, err
}

func (resp *Response) GetDom() (*goquery.Document, error) {
	if resp.dom != nil {
		return resp.dom, nil
	}
	text, err := resp.GetText()
	if err != nil {
		return nil, err
	}
	resp.dom, err = goquery.NewDocumentFromReader(bytes.NewReader(text))
	return resp.dom, err
}
