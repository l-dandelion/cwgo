package processor

import (
	"errors"

	"github.com/l-dandelion/cwgo/local/model"
	"github.com/l-dandelion/cwgo/local/processor/mysql-processor"
	"github.com/l-dandelion/cwgo/module"
)

/*
 * 根据规则生成处理器
 */
func GenItemProcessor(processor *model.Processor) (module.ProcessItem, error) {
	if processor == nil {
		return nil, errors.New("Nil processor.")
	}
	switch processor.Type {
	case model.PROCESSOR_TYPE_MYSQL:
		return mysql_processor.DefaultMysqlProcessor, nil
	}
	return nil, errors.New("Unsupport processor type.")
}

/*
 * 根据规则生成多个处理器
 */
func GenItemProcessors(processors []*model.Processor) ([]module.ProcessItem, error) {
	if len(processors) == 0 {
		return nil, errors.New("Empty processor list.")
	}
	result := []module.ProcessItem{}
	for _, processor := range processors {
		tmp, err := GenItemProcessor(processor)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
