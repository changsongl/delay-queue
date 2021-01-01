package config

import "gopkg.in/yaml.v2"

type DecodeFunc func([]byte, interface{}) error

type Decoder interface {
	GetYamlFunc() DecodeFunc
}

type decoder struct {
	yamlFunc DecodeFunc
}

func NewDecoder() Decoder {
	return &decoder{
		yamlFunc: yaml.Unmarshal,
	}
}

func (d *decoder) GetYamlFunc() DecodeFunc {
	return d.yamlFunc
}
