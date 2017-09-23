package tranhttp

import (
	"bitbucket.org/code_horse/pegasus/network"
)

type Options struct {
	network.Options
}

func NewOptions() *Options {
	opts := &Options{}
	opts.Fields = make(map[string]map[string]string)
	return opts
}

func (o *Options) BuildFromBytes(bt []byte) *Options {
	o.Options = *network.NewOptions().Unmarshal(bt)
	return o
}

func (o *Options) Build(options *network.Options) *Options {
	o.Options = *options
	return o
}

func (o *Options) GetOptions() *network.Options {
	return &o.Options
}

func (o *Options) SetPath(path string) *Options {
	o.Path = path
	return o
}

func (o *Options) GetPath(path string) string {
	return o.Path
}

func (o *Options) SetGetMethod() *Options {
	if o.Fields["METHOD"] == nil {
		o.Fields["METHOD"] = make(map[string]string)
	}
	o.Fields["METHOD"]["VALUE"] = "GET"
	return o
}

func (o *Options) SetPostMethod() *Options {
	if o.Fields["METHOD"] == nil {
		o.Fields["METHOD"] = make(map[string]string)
	}
	o.Fields["METHOD"]["VALUE"] = "POST"
	return o
}

func (o *Options) SetPutMethod() *Options {
	if o.Fields["METHOD"] == nil {
		o.Fields["METHOD"] = make(map[string]string)
	}
	o.Fields["METHOD"]["VALUE"] = "PUT"
	return o
}

func (o *Options) SetPatchMethod() *Options {
	if o.Fields["METHOD"] == nil {
		o.Fields["METHOD"] = make(map[string]string)
	}
	o.Fields["METHOD"]["VALUE"] = "PATCH"
	return o
}

func (o *Options) SetDeleteMethod() *Options {
	if o.Fields["METHOD"] == nil {
		o.Fields["METHOD"] = make(map[string]string)
	}
	o.Fields["METHOD"]["VALUE"] = "DELETE"
	return o
}

func (o *Options) GetMethod(path string) string {
	return o.Fields["METHOD"]["VALUE"]
}

func (o *Options) SetHeaders(headers map[string]string) *Options {
	if o.Fields["HEADER"] == nil {
		o.Fields["HEADER"] = make(map[string]string)
	}
	o.Fields["HEADER"] = headers
	return o
}

func (o *Options) SetHeader(key string, value string) *Options {
	if o.Fields["HEADER"] == nil {
		o.Fields["HEADER"] = make(map[string]string)
	}
	o.Fields["HEADER"][key] = value
	return o
}

func (o *Options) GetHeaders() map[string]string {
	return o.Fields["HEADER"]
}

func (o *Options) SetUrlParams(key string, value string) *Options {
	if o.Fields["URL-PARAMS"] == nil {
		o.Fields["URL-PARAMS"] = make(map[string]string)
	}
	o.Fields["URL-PARAMS"][key] = value
	return o
}

func (o *Options) GetUrlParams(value string) string {
	return o.Fields["URL-PARAMS"][value]
}

func (o *Options) SetQueryParams(value map[string]string) *Options {
	if o.Fields["QUERY-PARAMS"] == nil {
		o.Fields["QUERY-PARAMS"] = make(map[string]string)
	}
	o.Fields["QUERY-PARAMS"] = value
	return o
}

func (o *Options) SetQueryParam(key string, value string) *Options {
	if o.Fields["QUERY-PARAMS"] == nil {
		o.Fields["QUERY-PARAMS"] = make(map[string]string)
	}
	o.Fields["QUERY-PARAMS"][key] = value
	return o
}

func (o *Options) GetQueryParams(value string) string {
	return o.Fields["QUERY-PARAMS"][value]
}
