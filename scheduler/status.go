package scheduler

import (
	"errors"
	"fmt"
)

const (
	RUNNING_STATUS_UNPREPARED int8 = iota
	RUNNING_STATUS_PREPARING
	RUNNING_STATUS_PREPARED
	RUNNING_STATUS_STARTING
	RUNNING_STATUS_STARTED
	RUNNING_STATUS_PAUSING
	RUNNING_STATUS_PAUSED
	RUNNING_STATUS_STOPPING
	RUNNING_STATUS_STOPPED

	RUNNING_STATUS_UNPREPARED_DESC = "未初始化"
	RUNNING_STATUS_PREPARING_DESC  = "初始化中"
	RUNNING_STATUS_PREPARED_DESC   = "已初始化"
	RUNNING_STATUS_STARTING_DESC   = "启动中"
	RUNNING_STATUS_STARTED_DESC    = "已启动"
	RUNNING_STATUS_PAUSING_DESC    = "暂停中"
	RUNNING_STATUS_PAUSED_DESC     = "已暂停"
	RUNNING_STATUS_STOPPING_DESC   = "终止中"
	RUNNING_STATUS_STOPPED_DESC    = "已终止"
)

//状态描述
var RunningStatusDesc = map[int8]string{
	RUNNING_STATUS_UNPREPARED: RUNNING_STATUS_UNPREPARED_DESC,
	RUNNING_STATUS_PREPARING:  RUNNING_STATUS_PREPARING_DESC,
	RUNNING_STATUS_PREPARED:   RUNNING_STATUS_PREPARED_DESC,
	RUNNING_STATUS_STARTING:   RUNNING_STATUS_STARTING_DESC,
	RUNNING_STATUS_STARTED:    RUNNING_STATUS_STARTED_DESC,
	RUNNING_STATUS_PAUSING:    RUNNING_STATUS_PAUSING_DESC,
	RUNNING_STATUS_PAUSED:     RUNNING_STATUS_PAUSED_DESC,
	RUNNING_STATUS_STOPPING:   RUNNING_STATUS_STOPPING_DESC,
	RUNNING_STATUS_STOPPED:    RUNNING_STATUS_STOPPED_DESC,
}

/*
 * 检查是否可以进行该状态转换
 */
func checkStatus(currentStatus, wantedStatus int8) (err error) {
	switch currentStatus {
	case RUNNING_STATUS_PREPARING:
		err = errors.New("The scheduler is being initializd!")
	case RUNNING_STATUS_STARTING:
		err = errors.New("The scheduler is being started!")
	case RUNNING_STATUS_STOPPING:
		err = errors.New("The scheduler is being stopped!")
	case RUNNING_STATUS_PAUSING:
		err = errors.New("The scheduler is being paused!")
	}
	if err != nil {
		return
	}

	if currentStatus == RUNNING_STATUS_UNPREPARED &&
		wantedStatus != RUNNING_STATUS_PREPARING {
		err = errors.New("The scheduler has not yet been initialized!")
	}

	switch wantedStatus {
	case RUNNING_STATUS_PREPARING:
		switch currentStatus {
		case RUNNING_STATUS_STARTED:
			err = errors.New("the scheduler has been started!")
		case RUNNING_STATUS_PAUSED:
			err = errors.New("the scheduler has not been stopped!")
		}
	case RUNNING_STATUS_STARTING:
		switch currentStatus {
		case RUNNING_STATUS_UNPREPARED:
			err = errors.New("the scheduler has not been initialized!")
		case RUNNING_STATUS_STARTED:
			err = errors.New("the scheduler has been started!")
		}
	case RUNNING_STATUS_PAUSING:
		if currentStatus != RUNNING_STATUS_STARTED {
			err = errors.New("the scheduler has not been started!")
		}
	case RUNNING_STATUS_STOPPING:
		if currentStatus != RUNNING_STATUS_STARTED &&
			currentStatus != RUNNING_STATUS_PAUSED {
			err = errors.New("the scheduler has not been started!")
		}
	default:
		err = fmt.Errorf("unsupported wanted status for check! (wantedStatus: %d)",
			wantedStatus)
	}
	return
}

/*
 * 获取状态描述
 */
func GetStatusDescription(status int8) string {
	desc, ok := RunningStatusDesc[status]
	if !ok {
		return "Unknow"
	}
	return desc
}
