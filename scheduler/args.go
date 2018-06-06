package scheduler

import (
	"errors"

	"github.com/l-dandelion/cwgo/module"
)

type Args interface {
	Check() error
}

type RequestArgs struct {
	MaxDepth  int `json:"maxDepth"`
	MaxThread int `json:"maxThread"`
}

func (args *RequestArgs) Check() error {
	return nil
}

type ModuleArgs struct {
	Downloader module.Downloader
	Analyzer   module.Analyzer
	Pipeline   module.Pipeline
}

func (args ModuleArgs) Check() error {
	if args.Downloader == nil {
		return errors.New("Nil downloader.")
	}
	if args.Analyzer == nil {
		return errors.New("Nil analyzer.")
	}
	if args.Pipeline == nil {
		return errors.New("Nil pipeline.")
	}
	return nil
}
