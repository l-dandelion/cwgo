package module

import (
	"github.com/l-dandelion/cwgo/data"
)

type Counts struct {
	CalledCount    int64
	AcceptedCount  int64
	CompletedCount int64
	HandlingNumber int64
}

type SummaryStruct struct {
	Called    int64       `json:"called"`
	Accepted  int64       `json:"accepted"`
	Completed int64       `json:"completed"`
	Handling  int64       `json:"handling"`
	Extra     interface{} `json:"extra,omitempty"` // 额外信息
}

type Module interface {
	Counts() Counts
	Summary() SummaryStruct
	CalledCount() int64
	AcceptedCount() int64
	CompletedCount() int64
	HandlingNumber() int64
}

//下载器
type Downloader interface {
	Module
	Download(*data.Request) (*data.Response, error)
}

//分析器
type Analyzer interface {
	Module
	RespParsers() []ParseResponse
	Analyze(*data.Response) ([]data.Data, []error)
}

//响应解析器
type ParseResponse func(*data.Response) ([]data.Data, []error)

//条目处理管道
type Pipeline interface {
	Module
	ItemProcessors() []ProcessItem
	Send(item data.Item) []error

	FailFast() bool
	SetFailFast(failFast bool)
}

//条目处理器
type ProcessItem func(data.Item) (data.Item, error)
