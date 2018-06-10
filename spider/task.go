package spider

import (
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/local/model"
	"github.com/l-dandelion/cwgo/scheduler"
)

type Task struct {
	Id              string                `json:"id"`
	Name            string                `json:"name"`
	RequestArgs     scheduler.RequestArgs `json:"request_args"`
	RespParsers     []*model.Parser       `json:"resp_parsers"`
	ItemProcessors  []*model.Processor    `json:"item_processors"`
	InitialRequests []*data.Request        `json:"initial_requests"`
	FastFail        bool                  `json:"fast_fail"`
}
