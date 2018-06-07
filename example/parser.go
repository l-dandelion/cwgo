package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/utils"

)

// function for parsing response
func parseATag(resp *data.Response) ([]data.Data, []error) {
	reqURL := resp.HTTPResp().Request.URL
	httpResp := resp.HTTPResp()
	//TODO: 支持更多的HTTP响应状态。
	if httpResp.StatusCode != 200 {
		err := fmt.Errorf(
			fmt.Sprintf("Unsupported status code %d! (httpResponse: %v)",
				httpResp.StatusCode, httpResp))
		return nil, []error{err}
	}
	dom, err := resp.GetDom()
	if err != nil {
		return nil, []error{err}
	}
	dataList := []data.Data{}
	errList := []error{}
	dom.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || href == "#" || href == "/" {
			return
		}
		trimHref := strings.TrimSpace(href)
		lowHref := strings.ToLower(trimHref)
		fdStart := strings.Index(lowHref, "javascript")
		if fdStart == 0 {
			return
		}
		aURL, err := utils.ParseURL(lowHref)
		if err != nil {
			errList = append(errList, err)
			return
		}
		if !aURL.IsAbs() {
			aURL = reqURL.ResolveReference(aURL)
		}
		httpReq, err := http.NewRequest("GET", aURL.String(), nil)
		if err != nil {
			errList = append(errList, err)
			return
		}
		req := data.NewRequest(httpReq)
		dataList = append(dataList, req)
	})
	item := data.Item{
		"URL":     reqURL.String(),
		"DirPath": "result/",
		"Reader":  resp.HTTPResp().Body,
		"Etx":     ".html",
	}
	dataList = append(dataList, item)
	return dataList, errList
}

func parseATag2(resp *data.Response) ([]data.Data, []error) {
	matchedContentType := false
	httpResp := resp.HTTPResp()
	reqURL := httpResp.Request.URL
	if httpResp.Header != nil {
		contentTypes := httpResp.Header["Content-Type"]
		for _, contentType := range contentTypes {
			if strings.Index(contentType, "text/html") == 0 {
				matchedContentType = true
				break
			}
		}
	}
	dataList := []data.Data{}
	errList := []error{}
	if matchedContentType {
		dom, err := resp.GetDom()
		if err != nil {
			return nil, []error{err}
		}
		dom.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists || href == "" || href == "#" || href == "/" {
				return
			}
			trimHref := strings.TrimSpace(href)
			lowHref := strings.ToLower(trimHref)
			fdStart := strings.Index(lowHref, "javascript")
			if fdStart == 0 {
				return
			}
			aURL, err := utils.ParseURL(lowHref)
			if err != nil {
				errList = append(errList, err)
				return
			}
			if !aURL.IsAbs() {
				aURL = reqURL.ResolveReference(aURL)
			}
			httpReq, err := http.NewRequest("GET", aURL.String(), nil)
			if err != nil {
				errList = append(errList, err)
				return
			}
			req := data.NewRequest(httpReq)
			dataList = append(dataList, req)
		})
		dom.Find("img").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("src")
			if !exists || href == "" || href == "#" || href == "/" {
				return
			}
			trimHref := strings.TrimSpace(href)
			lowHref := strings.ToLower(trimHref)
			fdStart := strings.Index(lowHref, "javascript")
			if fdStart == 0 {
				return
			}
			aURL, err := utils.ParseURL(lowHref)
			if err != nil {
				errList = append(errList, err)
				return
			}
			if !aURL.IsAbs() {
				aURL = reqURL.ResolveReference(aURL)
			}
			httpReq, err := http.NewRequest("GET", aURL.String(), nil)
			if err != nil {
				errList = append(errList, err)
				return
			}
			req := data.NewRequest(httpReq)
			dataList = append(dataList, req)
		})

		item := data.Item{
			"URL":     reqURL.String(),
			"DirPath": "temp/500px/htmls",
			"Reader":  resp.HTTPResp().Body,
			"Etx":     ".html",
		}
		dataList = append(dataList, item)
	}
	return dataList, errList
}

func parseImgTag(resp *data.Response) ([]data.Data, []error) {
	pictureFormat := ""
	httpResp := resp.HTTPResp()
	if httpResp.Header != nil {
		contentTypes := httpResp.Header["Content-Type"]
		contentType := ""
		for _, ct := range contentTypes {
			if strings.Index(ct, "image") == 0 {
				contentType = ct
				break
			}
		}

		index1 := strings.Index(contentType, "/")
		index2 := strings.Index(contentType, ";")
		if index1 > 0 {
			if index2 < 0 {
				pictureFormat = contentType[index1+1:]
			} else if index1 < index2 {
				pictureFormat = contentType[index1+1 : index2]
			}
		}
	}
	dataList := []data.Data{}
	errList := []error{}
	if pictureFormat != "" {
		reqURL := resp.HTTPResp().Request.URL
		item := data.Item{
			"URL":     reqURL.String(),
			"DirPath": "temp/500px/imgs/",
			"Reader":  resp.HTTPResp().Body,
			"Etx":     "." + pictureFormat,
		}
		dataList = append(dataList, item)
	}
	return dataList, errList
}
