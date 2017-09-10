package network

import (
	"encoding/json"
)

type Options struct {
	Path    string
	Fields  map[string]map[string]string
}

func NewOptions() *Options {
	return &Options{Fields: make(map[string]map[string]string)}
}

func (c *Options) Marshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Options) Unmarshal(data []byte) *Options {
	err := json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Options) CreateNewField(key string) {
	c.Fields[key] = make(map[string]string)
}
