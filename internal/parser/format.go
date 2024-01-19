package parser

import "gopkg.in/yaml.v3"

type DecodeFunc func([]byte, interface{}) error

type Format struct {
	Begin string
	End   string

	DecodeFn DecodeFunc
}

func NewYamlFrontMatterFormat() *Format {
	return &Format{
		Begin:    "---",
		End:      "---",
		DecodeFn: yaml.Unmarshal,
	}
}
