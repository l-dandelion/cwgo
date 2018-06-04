package analyzer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/lib/reader"
	"github.com/l-dandelion/cwgo/module"
	"github.com/pkg/errors"
)

type myAnalyzer struct {
	module.ModuleInternal
	respParsers []module.ParseResponse
}

//获取解析函数
func (analyzer *myAnalyzer) RespParsers() []module.ParseResponse {
	parsers := make([]module.ParseResponse, len(analyzer.respParsers))
	copy(parsers, analyzer.respParsers)
	return parsers
}

//解析响应获取数据
func (analyzer *myAnalyzer) Analyze(
	resp *data.Response) (dataList []data.Data, errList []error) {
	analyzer.IncrHandlingNumber()
	defer analyzer.DecrHandlingNumber()
	analyzer.IncrCalledCount()
	errList = []error{}

	// 检查参数
	if resp == nil {
		err := errors.New("Nil response.")
		errList = append(errList, err)
		return
	}
	httpResp := resp.HTTPResp()
	if httpResp == nil {
		err := errors.New("Nil HTTP response.")
		errList = append(errList, err)
		return
	}
	req := resp.Request()
	if req == nil {
		err := errors.New("Nil request.")
		errList = append(errList, err)
		return
	}
	httpReq := req.HTTPReq()
	if httpReq == nil {
		err := errors.New("Nil HTTP request.")
		errList = append(errList, err)
		return
	}
	reqURL := httpReq.URL
	if reqURL == nil {
		err := errors.New("Nil request URL.")
		errList = append(errList, err)
		return
	}
	analyzer.IncrAcceptedCount()

	respDepth := resp.Depth()
	log.Infof("Parse the response (URL: %s, depth: %d)... \n", reqURL, respDepth)
	if httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	//重复读body
	multiReader, err := reader.NewMultipleReader(httpResp.Body)
	if err != nil {
		errList = append(errList, err)
		return
	}
	dataList = []data.Data{}
	for _, respParser := range analyzer.respParsers {
		if httpResp.Body != nil {
			httpResp.Body.Close()
		}
		httpResp.Body = multiReader.Reader()
		pDataList, pErrList := respParser(resp)
		if pDataList != nil {
			for _, mdata := range pDataList {
				dataList = appendDataList(dataList, mdata, respDepth)
			}
		}
		if pErrList != nil {
			for _, err := range pErrList {
				errList = append(errList, err)
			}
		}
	}
	if len(errList) == 0 {
		analyzer.IncrCompletedCount()
	}
	return
}

func appendDataList(dataList []data.Data, mdata data.Data, respDepth int) []data.Data {
	if mdata == nil {
		return dataList
	}
	req, ok := mdata.(*data.Request)
	if ok {
		req.SetDepth(respDepth + 1)
	}
	return append(dataList, mdata)
}