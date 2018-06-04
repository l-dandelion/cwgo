package downloader

import (
	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/cwgo/data"
	"github.com/cwgo/module"
)

//重试次数
var RetryTimes = 3

//只能通过new创建downloader
func New(client *http.Client) (module.Downloader, error) {
	if client == nil {
		return nil, errors.New("Nil client.")
	}
	return &myDownloader{
		ModuleInternal: module.NewModuleInternal(),
		httpClient:     client,
	}, nil
}

type myDownloader struct {
	module.ModuleInternal
	httpClient *http.Client
}

/*
 * 根据请求下载并返回响应
 */
func (downloader *myDownloader) Download(req *data.Request) (*data.Response, error) {
	downloader.IncrHandlingNumber()
	defer downloader.DecrHandlingNumber()
	downloader.IncrCalledCount()

	//检查请求参数
	if req == nil {
		return nil, errors.New("Nil request.")
	}
	if req.HTTPReq() == nil {
		return nil, errors.New("Nil HTTP request.")
	}
	downloader.IncrAcceptedCount()

	var (
		httpResp *http.Response
		err      error
	)

	log.Infof("Do the request (URL: %s, depth: %d)... \n",
		req.HTTPReq().URL, req.Depth())

	//尝试下载
	for i := 0; i < RetryTimes; i++ {
		httpResp, err = downloader.httpClient.Do(req.HTTPReq())
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}
	resp := data.NewResponse(req, httpResp)
	downloader.IncrCompletedCount()
	return resp, nil
}
