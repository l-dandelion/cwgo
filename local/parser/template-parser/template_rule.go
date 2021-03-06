package template_parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/lib/parseurl"
	"github.com/l-dandelion/cwgo/local/model"
	"github.com/l-dandelion/cwgo/utils"
)

func TemplateRuleProcess(parser *model.Parser, resp *data.Response) (dataList []data.Data, errorList []error) {
	dataList = []data.Data{}
	errorList = []error{}
	rule := parser.TemplateRule.Rule

	doc, err := resp.GetDom()
	if err != nil {
		errorList = append(errorList, err)
		return
	}

	if len(parser.RegUrl) > 0 {
		doc.Find("a").Each(func(i int, sel *goquery.Selection) {
			href, _ := sel.Attr("href")
			href, err = utils.GetComplateUrl(resp.HTTPRequest().URL, href)
			if err != nil {
				errorList = append(errorList, err)
				return
			}

			httpReq, err := http.NewRequest("GET", href, nil)
			if err != nil {
				errorList = append(errorList, err)
				return
			}
			dataList = append(dataList, data.NewRequest(httpReq))

		})
	}

	resultType := "map"
	rootSel := ""

	v, ok := rule["node"]
	if ok {
		contentInfo := strings.Split(v, "|")
		resultType = contentInfo[0]
		rootSel = contentInfo[1]
	}

	if resultType == "array" {
		doc.Find(rootSel).Each(func(i int, s *goquery.Selection) {
			mdata := getMapFromDom(rule, s)
			if mdata == nil {
				return
			}
			dataList = append(dataList, data.Item(mdata))
			if len(parser.AddQueue) > 0 {
				urls := parseurl.ParseReqUrl(parser.AddQueue, mdata)
				fmt.Println(urls, " ", parser.AddQueue, " ", mdata)
				for _, u := range urls {
					httpReq, err := http.NewRequest("GET", u, nil)
					if err != nil {
						errorList = append(errorList, err)
						return
					}
					dataList = append(dataList, data.NewRequest(httpReq))
				}
			}

		})
	}

	if resultType == "map" {
		mdata := getMapFromDom(rule, doc.Selection)
		dataList = append(dataList, data.Item(mdata))
		if len(parser.AddQueue) > 0 {
			urls := parseurl.ParseReqUrl(parser.AddQueue, mdata)
			for _, u := range urls {
				httpReq, err := http.NewRequest("GET", u, nil)
				if err != nil {
					errorList = append(errorList, err)
					return
				}
				dataList = append(dataList, data.NewRequest(httpReq))
			}
		}
	}

	return
}

func getMapFromDom(rule map[string]string, node *goquery.Selection) map[string]interface{} {

	result := make(map[string]interface{})

	isNull := true

	for key, value := range rule {

		if key == "node" {
			continue
		}

		rules := strings.Split(value, "|")
		ValueType := strings.Split(rules[0], ".")

		if len(rules) < 2 {
			result[key] = value
			continue
		}

		s := node.Find(rules[1])
		switch ValueType[0] {
		case "text":
			result[key] = s.Text()
		case "html":
			result[key], _ = s.Html()
		case "attr":
			if len(ValueType) < 2 {
				continue
			}
			result[key], _ = s.Attr(ValueType[1])
		case "texts":
			arr := []string{}
			s.Each(func(i int, sel *goquery.Selection) {
				text := sel.Text()
				arr = append(arr, text)
			})
			j, _ := json.Marshal(arr)
			result[key] = string(j)
		case "htmls":
			arr := []string{}
			s.Each(func(i int, sel *goquery.Selection) {
				html, _ := s.Html()
				arr = append(arr, html)
			})
			j, _ := json.Marshal(arr)
			result[key] = string(j)
		case "attrs":
			arr := []string{}
			attr := ""
			s.Each(func(i int, sel *goquery.Selection) {
				if len(ValueType) >= 2 {
					attr, _ = sel.Attr(ValueType[1])
					arr = append(arr, attr)
				}
			})
			result[key] = arr
		default:
			result[key] = value
		}
		res, ok := result[key].(string)
		if ok || len(res) != 0 {
			isNull = false
		}
	}

	if isNull == true {
		return nil
	}

	return result
}
