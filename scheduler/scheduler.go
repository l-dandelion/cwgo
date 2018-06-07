package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/lib/buffer"
	"github.com/l-dandelion/cwgo/lib/cmap"
	"github.com/l-dandelion/cwgo/lib/pool"
	"github.com/l-dandelion/cwgo/module"
	"github.com/pkg/errors"
)

//缓冲器大小
var (
	BufferCap       uint32 = 1000
	MaxBufferNumber uint32 = 10000
)

type Scheduler interface {
	Name() string
	Init(requestArgs RequestArgs, moduleArgs ModuleArgs) error
	Start(initialReqs []*data.Request) error
	Pause() error
	Recover() error
	Stop() error
	Status() int8
	ErrorChan() <-chan error
	Idle() bool
}

type myScheduler struct {
	name            string
	maxDepth        int
	reqBufferPool   buffer.Pool
	respBufferPool  buffer.Pool
	itemBufferPool  buffer.Pool
	errorBufferPool buffer.Pool
	urlMap          cmap.ConcurrentMap
	ctx             context.Context
	cancelFunc      context.CancelFunc
	status          int8
	statusLock      sync.RWMutex
	downloader      module.Downloader
	analyzer        module.Analyzer
	pipeline        module.Pipeline
	summary         SchedSummary
	downloaderPool  pool.Pool
	analyzerPool    pool.Pool
	pipelinePool    pool.Pool
}

/*
 * create an instance of interface Scheduler by name
 */
func New(name string) Scheduler {
	return &myScheduler{name: name}
}

/*Get*/

func (sched *myScheduler) Name() string {
	return sched.name
}

func (sched *myScheduler) Init(
	requestArgs RequestArgs,
	moduleArgs ModuleArgs) (err error) {
	//检查状态
	log.Info("Check status for initialization...")
	oldStatus, err := sched.checkAndSetStatus(RUNNING_STATUS_PREPARING)
	if err != nil {
		return
	}
	//检查是否初始化成功，不成功还原状态
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = RUNNING_STATUS_PREPARED
		}
		sched.statusLock.Unlock()
	}()

	log.Info("Check request arguments...")
	if err = requestArgs.Check(); err != nil {
		return
	}
	log.Info("Request arguments are valid.")

	log.Info("Check module arguments...")
	if err = moduleArgs.Check(); err != nil {
		return
	}
	log.Info("Module arguments are valid.")

	log.Info("Initialize Scheduler's fields...")
	sched.maxDepth = requestArgs.MaxDepth
	log.Infof("-- Max depth: %d", sched.maxDepth)

	sched.urlMap, _ = cmap.NewConcurrentMap(16, nil)
	log.Infof("-- URL map: length: %d, concurrency: %d",
		sched.urlMap.Len(), sched.urlMap.Concurrency())
	sched.downloader = moduleArgs.Downloader
	sched.analyzer = moduleArgs.Analyzer
	sched.pipeline = moduleArgs.Pipeline
	//控制并发量
	sched.downloaderPool = pool.New(requestArgs.MaxThread)
	sched.analyzerPool = pool.New(requestArgs.MaxThread)
	sched.pipelinePool = pool.New(requestArgs.MaxThread)
	sched.initBufferPool()
	sched.resetContext()
	sched.summary = newSchedSummary(requestArgs, moduleArgs, sched)

	log.Info("Scheduler has been initialized.")
	return
}

/*
 * 开启
 */
func (sched *myScheduler) Start(initialReqs []*data.Request) (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal Scheduler error: %s", p)
			log.Fatal(errMsg)
			err = errors.New(errMsg)
		}
	}()
	log.Info("Start Scheduler ...")
	log.Info("Check status for start ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(RUNNING_STATUS_STARTING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = RUNNING_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()

	sched.download()
	sched.analyze()
	sched.pick()
	log.Info("The Scheduler has been started.")
	for _, req := range initialReqs {
		sched.sendReq(req)
	}
	return nil
}

/*
 * 暂停
 */
func (sched *myScheduler) Pause() (err error) {
	//check status
	log.Info("Pause Scheduler ...")
	log.Info("Check status for pause ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(RUNNING_STATUS_PAUSING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = RUNNING_STATUS_PAUSED
		}
		sched.statusLock.Unlock()
	}()
	log.Info("Scheduler has been paused.")
	return nil
}

/*
 * 恢复
 */
func (sched *myScheduler) Recover() (err error) {
	log.Info("Recover Scheduler ...")
	log.Info("Check status for recover ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(RUNNING_STATUS_STARTING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = RUNNING_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	log.Info("Scheduler has been recovered.")
	return nil
}

/*
 * 终止
 */
func (sched *myScheduler) Stop() (err error) {
	log.Info("Stop Scheduler ...")
	log.Info("Check status for stop ...")
	var oldStatus int8
	oldStatus, err = sched.checkAndSetStatus(RUNNING_STATUS_STOPPING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = RUNNING_STATUS_STOPPED
		}
		sched.statusLock.Unlock()
	}()

	sched.cancelFunc()
	sched.reqBufferPool.Close()
	sched.respBufferPool.Close()
	sched.itemBufferPool.Close()
	sched.errorBufferPool.Close()
	log.Info("Scheduler has been stopped.")
	return nil
}

/*
 * 初始化缓冲池
 */
func (sched *myScheduler) initBufferPool() {
	sched.reqBufferPool, _ = buffer.NewPool(BufferCap, MaxBufferNumber)
	sched.respBufferPool, _ = buffer.NewPool(BufferCap, MaxBufferNumber)
	sched.itemBufferPool, _ = buffer.NewPool(BufferCap, MaxBufferNumber)
	sched.errorBufferPool, _ = buffer.NewPool(BufferCap, MaxBufferNumber)
}

/*
 * 检查并进行状态设置
 */
func (sched *myScheduler) checkAndSetStatus(wantedStatus int8) (oldStatus int8, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}

/*
 * 重置
 */
func (sched *myScheduler) resetContext() {
	sched.ctx, sched.cancelFunc = context.WithCancel(context.Background())
}

/*
 * 检查是否终止
 */
func (sched *myScheduler) canceled() bool {
	select {
	case <-sched.ctx.Done():
		return true
	default:
		return false
	}
}

/*
 * 获取运行状态
 */
func (sched *myScheduler) Status() int8 {
	sched.statusLock.RLock()
	defer sched.statusLock.RUnlock()
	return sched.status
}

/*
 * 是否完成
 */
func (sched *myScheduler) Idle() bool {
	if sched.downloader.HandlingNumber() > 0 ||
		sched.analyzer.HandlingNumber() > 0 ||
		sched.pipeline.HandlingNumber() > 0 {
		return false
	}
	if sched.reqBufferPool.Total() > 0 ||
		sched.respBufferPool.Total() > 0 ||
		sched.itemBufferPool.Total() > 0 {
		return false
	}
	return true
}

/*
 * 获取错误通道
 */
func (sched *myScheduler) ErrorChan() <-chan error {
	errBuffer := sched.errorBufferPool
	errCh := make(chan error, errBuffer.BufferCap())
	go func(errBuffer buffer.Pool, errCh chan error) {
		for {
			//stopped
			if sched.canceled() {
				close(errCh)
				break
			}
			//paused
			if sched.Status() == RUNNING_STATUS_PAUSED {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			datum, err := errBuffer.Get()
			if err != nil {
				log.Warnln("The error buffer pool was closed. Break error reception.")
				close(errCh)
				break
			}
			err, ok := datum.(error)
			if !ok {
				err = fmt.Errorf("Incorrect error type: %T", datum)
				sched.sendError(err)
				continue
			}
			if sched.canceled() {
				close(errCh)
				break
			}
			errCh <- err
		}
	}(errBuffer, errCh)
	return errCh
}
