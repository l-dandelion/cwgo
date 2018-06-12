package crawler

import "errors"

var (
	ERR_SPIDER_NOT_FOUND     = errors.New("Spider not found.")
	ERR_SPIDER_NAME_REPEATED = errors.New("Spider name repeated.")
)
