package trangrpc

import (
	"bitbucket.org/code_horse/pegasus/network"
)

type Options struct {
	network.Options
}

func NewOptions() *Options {
	options := &Options{}
	options.Fields = make(map[string]map[string]string)
	return options
}

func BuildOptions(networkOptions *network.Options) *Options {
	return &Options{Options: *networkOptions}
}

func BuildOptionsFromBytes(data []byte) *Options {
	return &Options{Options: *network.NewOptions().Unmarshal(data)}
}

func (o *Options) GetOptions() *network.Options {
	return &o.Options
}

func (o *Options) initMapper(key string) {
	if o.Fields[key] == nil {
		o.Fields[key] = make(map[string]string)
	}
}
