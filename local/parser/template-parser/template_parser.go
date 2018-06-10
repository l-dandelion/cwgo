package template_parser

import (
	"github.com/l-dandelion/cwgo/data"
	"github.com/l-dandelion/cwgo/module"
	"github.com/l-dandelion/cwgo/local/model"
	"github.com/l-dandelion/cwgo/local/parser/filter"
)

func GenTemplateParser(parser *model.Parser) module.ParseResponse {
	return func(resp *data.Response) ([]data.Data, []error) {
		if len(parser.RegUrl) > 0 && !filter.Filter(resp.HTTPRequest().URL.String(), parser.RegUrl) {
			return nil, nil
		}
		return TemplateRuleProcess(parser, resp)
	}
}
