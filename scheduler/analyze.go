package scheduler

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/l-dandelion/cwgo/data"
)

/*
 * 开启解析goroutine
 */
func (sched *myScheduler) analyze() {
	go func() {
		for {
			//stopped
			if sched.canceled() {
				break
			}
			//paused
			if sched.Status() == RUNNING_STATUS_PAUSED {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			datum, err := sched.respBufferPool.Get()
			if err != nil {
				log.Warnln("The response buffer pool was closed. Break response reception.")
				break
			}
			sched.analyzerPool.Add()
			go func(datum interface{}) {
				defer sched.analyzerPool.Done()
				resp, ok := datum.(*data.Response)
				if !ok {
					err := fmt.Errorf("Incorrect response type: %T", datum)
					sched.sendError(err)
					return
				}
				sched.analyzeOne(resp)
			}(datum)
		}
		sched.analyzerPool.Wait()
	}()
}

/*
 * 处理一个响应
 */
func (sched *myScheduler) analyzeOne(resp *data.Response) {
	if resp == nil {
		return
	}
	if sched.canceled() {
		return
	}
	analyzer := sched.analyzer
	dataList, errs := analyzer.Analyze(resp)
	if dataList != nil {
		for _, mdata := range dataList {
			if mdata == nil {
				continue
			}
			switch d := mdata.(type) {
			case *data.Request:
				d.SetDepth(resp.Depth() + 1)
				sched.sendReq(d)
			case data.Item:
				sched.sendItem(d)
			default:
				err := fmt.Errorf("Unsupported data type: %T (data: %#v)", mdata, mdata)
				sched.sendError(err)
			}
		}
	}
	if errs != nil {
		for _, err := range errs {
			sched.sendError(err)
		}
	}
}
