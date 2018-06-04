package pipeline

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/cwgo/data"
	"github.com/cwgo/module"
)

type myPipeline struct {
	module.ModuleInternal
	itemProcessors []module.ProcessItem
	//快速失败，当为真的时候，如果在某个处理器出错，后续处理器自动跳过
	failFast bool
}

/*Get*/

/*
 * 获取条目处理器列表
 */
func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	processors := make([]module.ProcessItem, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
}

/*Set*/

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

/*Other*/

/*
 * 将条目依次传入处理器中进行处理
 */
func (pipeline *myPipeline) Send(item data.Item) (errs []error) {
	pipeline.IncrHandlingNumber()
	defer pipeline.DecrHandlingNumber()
	pipeline.IncrCalledCount()
	if item == nil {
		errs = append(errs, errors.New("Nil item."))
		return
	}
	pipeline.IncrAcceptedCount()
	log.Infof("Process item %+v... \n", item)
	currentItem := item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			if pipeline.failFast {
				//快速失败
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	if len(errs) == 0 {
		pipeline.IncrCompletedCount()
	}
	return
}



