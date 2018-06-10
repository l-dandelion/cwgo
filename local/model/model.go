package model

const (
	PARSER_TYPE_TEMPLATE = "template"
)

type TemplateRule struct {
	Rule map[string]string
}

type JsonRule struct {
	Rule map[string]string
}

type Parser struct {
	RegUrl       []string     `json:"reg_url"`
	Type         string       `json:"type"`
	TemplateRule TemplateRule `json:"template_rule"`
	JsonRule     JsonRule     `json:"json_rule"`
	AddQueue     []string     `json:"add_queue"`
}


const (
	PROCESSOR_TYPE_MYSQL = "mysql"
)

type Processor struct {
	Type string
	Rule map[string]string
}