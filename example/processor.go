package main

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/l-dandelion/cwgo/utils"
	"github.com/l-dandelion/cwgo/data"

)

func process(item data.Item) (data.Item, error) {
	dirPath := item["DirPath"].(string)
	murl := item["URL"].(string)
	h := md5.New()
	h.Write([]byte(murl))
	etx := item["Etx"].(string)
	fileName := fmt.Sprintf("%x", h.Sum(nil)) + etx
	reader := item["Reader"].(io.Reader)
	err := utils.SaveFileByReader(dirPath, fileName, reader)
	if err != nil {
		return nil, err
	}
	return nil, nil
}