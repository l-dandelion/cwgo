package spider

import (
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/local/parser"
	"github.com/l-dandelion/cwgo/local/processor"
	"github.com/l-dandelion/cwgo/module/analyzer"
	"github.com/l-dandelion/cwgo/module/downloader"
	"github.com/l-dandelion/cwgo/module/pipeline"
	"github.com/l-dandelion/cwgo/scheduler"
)

type Spider interface {
	Init() error
	Start() error
	Stop() error
	Pause() error
	Recover() error
}

type mySpider struct {
	Id              string
	Name            string
	RequestArgs     scheduler.RequestArgs
	ModuleArgs      scheduler.ModuleArgs
	InitialRequests []*data.Request
	sched           scheduler.Scheduler
}

/*
 * 根据任务生成爬虫
 */
func NewSpiderWithTask(task Task) (Spider, error) {
	parsers, err := parser.GenRespParsers(task.RespParsers)
	if err != nil {
		return nil, err
	}
	processors, err := processor.GenItemProcessors(task.ItemProcessors)
	if err != nil {
		return nil, err
	}
	d := downloader.New(GenDefaultClient())
	a := analyzer.New(parsers)
	p := pipeline.New(processors, task.FastFail)
	return &mySpider{
		Id:          task.Id,
		Name:        task.Name,
		RequestArgs: task.RequestArgs,
		ModuleArgs: scheduler.ModuleArgs{
			Downloader: d,
			Analyzer:   a,
			Pipeline:   p,
		},
		InitialRequests: task.InitialRequests,
		sched:           scheduler.New(task.Name),
	}, nil
}

/*
 * 初始化调度器
 */
func (spider *mySpider) Init() error {
	return spider.sched.Init(spider.RequestArgs, spider.ModuleArgs)
}

/*
 * 开启调度器
 */
func (spider *mySpider) Start() error {
	return spider.sched.Start(spider.InitialRequests)
}

/*
 * 终止调度器
 */
func (spider *mySpider) Stop() error {
	return spider.sched.Stop()
}

/*
 * 暂停调度器
 */
func (spider *mySpider) Pause() error {
	return spider.sched.Pause()
}

/*
 * 回复调度器
 */
func (spider *mySpider) Recover() error {
	return spider.sched.Recover()
}
