package scheduler

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/l-dandelion/cwgo/data"
)

/*
 * 开启下载goroutine
 */
func (sched *myScheduler) download() {
	go func() {
		for {
			// stopped
			if sched.canceled() {
				break
			}
			// paused
			if sched.Status() == RUNNING_STATUS_PAUSED {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			datum, err := sched.reqBufferPool.Get()
			if err != nil {
				log.Warnln("The request buffer pool was closed. Break request reception.")
				return
			}

			sched.downloaderPool.Add()
			go func(datum interface{}) {
				defer sched.downloaderPool.Done()
				req, ok := datum.(*data.Request)
				if !ok {
					err := fmt.Errorf("Incorrect request type: %T", datum)
					sched.sendError(err)
					return
				}
				sched.downloadOne(req)
			}(datum)
		}
		sched.downloaderPool.Wait()
	}()
}

/*
 * 下载一个请求
 */
func (sched *myScheduler) downloadOne(req *data.Request) {
	if req == nil {
		return
	}
	if sched.canceled() {
		return
	}
	downloader := sched.downloader
	resp, err := downloader.Download(req)
	if resp != nil {
		sched.sendResp(resp)
	}
	if err != nil {
		sched.sendError(err)
	}
}
