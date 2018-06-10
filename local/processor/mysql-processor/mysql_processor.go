package mysql_processor

import (
	"errors"

	"github.com/l-dandelion/cwgo/data"
)

var DefaultMysqlProcessor = func(item data.Item) (result data.Item, err error) {
	if _, ok := item["kind"]; !ok {
		return nil, errors.New("Kind is empty.")
	}
	kind := item["kind"].(string)
	if kind != "mysql" {
		return nil, nil
	}
	if _, ok := item["tableName"]; !ok {
		return nil, errors.New("Table name not found.")
	}
	tableName := item["tableName"].(string)
	delete(item, "tableName")
	delete(item, "kind")
	dbModel := NewDBModel(tableName, item)
	item["tableName"] = tableName
	item["kind"] = kind
	AddPrepare(dbModel)
	return nil, nil
}
