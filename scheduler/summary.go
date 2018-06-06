package scheduler

import (
	"encoding/json"

	"github.com/l-dandelion/cwgo/module"
)

type SchedSummary interface {
	Struct() SummaryStruct
	String() (string, error)
}

type SummaryStruct struct {
	Status     string               `json:"status"`
	Downloader module.SummaryStruct `json:"downloader"`
	Analyzer   module.SummaryStruct `json:"analyzer"`
	Pipeline   module.SummaryStruct `json:"pipeline"`

	//各个缓冲池的大小
	ReqBufferSize   uint64 `json:"reqBufferSize"`
	RespBufferSize  uint64 `json:"respBufferSize"`
	ItemBufferSize  uint64 `json:"itemBufferSize"`
	ErrorBufferSize uint64 `json:"errorBufferSize"`

	//url总数
	NumURL uint64 `json:"numUrl"`
}

func (ss *mySchedSummary) Struct() SummaryStruct {
	return SummaryStruct{
		Status:          GetStatusDescription(ss.sched.Status()),
		Downloader:      ss.sched.downloader.Summary(),
		Analyzer:        ss.sched.analyzer.Summary(),
		Pipeline:        ss.sched.pipeline.Summary(),
		ReqBufferSize:   ss.sched.reqBufferPool.Total(),
		RespBufferSize:  ss.sched.respBufferPool.Total(),
		ItemBufferSize:  ss.sched.itemBufferPool.Total(),
		ErrorBufferSize: ss.sched.errorBufferPool.Total(),
		NumURL:          ss.sched.urlMap.Len(),
	}
}

func (ss *mySchedSummary) String() (string, error) {
	b, err := json.MarshalIndent(ss.Struct(), "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

/*
 * create an instance of SchedSummary
 */
func newSchedSummary(requestArgs RequestArgs, moduleArgs ModuleArgs, sched *myScheduler) SchedSummary {
	if sched == nil {
		return nil
	}
	return &mySchedSummary{
		requestArgs: requestArgs,
		moduleArgs:  moduleArgs,
		maxDepth:    requestArgs.MaxDepth,
		sched:       sched,
	}
}

type mySchedSummary struct {
	requestArgs RequestArgs
	moduleArgs  ModuleArgs
	maxDepth    int
	sched       *myScheduler
}
