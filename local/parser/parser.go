package parser

import (
	"errors"

	"github.com/l-dandelion/cwgo/local/model"
	"github.com/l-dandelion/cwgo/local/parser/template-parser"
	"github.com/l-dandelion/cwgo/module"
)

/*
 * 根据规则生成解析器
 */
func GenRespParser(parser *model.Parser) (module.ParseResponse, error) {
	if parser == nil {
		return nil, errors.New("Nil parser.")
	}
	switch parser.Type {
	case model.PARSER_TYPE_TEMPLATE:
		return template_parser.GenTemplateParser(parser), nil
	}
	return nil, errors.New("Unsupport parser type.")
}

/*
 * 根据规则生成多个解析器
 */
func GenRespParsers(parsers []*model.Parser) ([]module.ParseResponse, error) {
	if len(parsers) == 0 {
		return nil, errors.New("Empty parser list.")
	}
	result := []module.ParseResponse{}
	for _, parser := range parsers {
		tmp, err := GenRespParser(parser)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}
	return result, nil
}
