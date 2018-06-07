package main

import (
	"time"

	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/module"
	"github.com/l-dandelion/cwgo/module/analyzer"
	"github.com/l-dandelion/cwgo/module/downloader"
	"github.com/l-dandelion/cwgo/module/pipeline"
	"github.com/l-dandelion/cwgo/scheduler"
	"net/http"
)

func main() {
	sched := scheduler.New("test")
	requestArgs := scheduler.RequestArgs{
		MaxThread: 10,
		MaxDepth:  3,
	}
	d := downloader.New(genHTTPClient())
	a := analyzer.New([]module.ParseResponse{parseATag2, parseImgTag})
	p := pipeline.New([]module.ProcessItem{process}, false)
	moduleArgs := scheduler.ModuleArgs{
		Downloader: d,
		Analyzer:   a,
		Pipeline:   p,
	}
	sched.Init(requestArgs, moduleArgs)
	httpReq, _ := http.NewRequest("GET", "http://pixabay.com", nil)
	req := data.NewRequest(httpReq)
	sched.Start([]*data.Request{req})
	time.Sleep(1000 * time.Second)
}
