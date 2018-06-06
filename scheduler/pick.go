package scheduler

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/l-dandelion/cwgo/data"
	"fmt"
)

/*
 * 开始收集goroutine
 */
func (sched *myScheduler) pick() {
	go func() {
		for {
			//stopped
			if sched.canceled() {
				break
			}
			//paused
			if sched.Status() == RUNNING_STATUS_PAUSING {
				time.Sleep(100*time.Millisecond)
				continue
			}
			datum, err := sched.itemBufferPool.Get()
			if err != nil {
				log.Warnln("The item buffer pool was closed. Break item reception.")
				break
			}
			sched.pipelinePool.Add()
			go func(datum interface{}){
				defer sched.pipelinePool.Done()
				item, ok := datum.(data.Item)
				if !ok {
					err := fmt.Errorf("Incorrect item type: %T, item=%+v", item, item)
					sched.sendError(err)
					return
				}
				sched.pickOne(item)
			}(datum)

		}
	}()
}

/*
 * 收集一个条目
 */
func (sched *myScheduler) pickOne(item data.Item) {
	if item == nil {
		return
	}
	if sched.canceled() {
		return
	}
	pipeline := sched.pipeline
	errs := pipeline.Send(item)
	if errs != nil {
		for _, err := range errs {
			sched.sendError(err)
		}
	}
}