package network

import (
	"encoding/json"
)

// Options struct defines the parameters that we can pass through the  network. It contains three layers of parameters
// that could be defined.
//
// The Headers
//  Header: Fields["HEADERS"]["ANYTHING"] = "STRING VALUE"
// The headers are used to pass options in order to tell the server to change a process or to check something.
//
// The Params
//  Header: Fields["PARAMS"]["ANYTHING"] = "STRING VALUE"
// The params are set by the server only and define the parameters that can be set. e.g. HTTP url or path params.
//
// Custom
//  Header: Fields["ANYTHING"]["ANYTHING"] = "STRING VALUE"
// Custom fields are used if we want to set something completely custom, it should  be avoided.
type Options struct {

	// Fields are general properties may client need to know
	Fields map[string]map[string]string
}

// NewOptions creates an Option object, initialize the struct properties and returns it.
func NewOptions() *Options {
	return &Options{Fields: make(map[string]map[string]string)}
}

// BuildOptions receives bytes as parameter and converts them to Option object.
func BuildOptions(data []byte) *Options {
	return NewOptions().Unmarshal(data)
}

// SetParams re-sets the hole parameters fields.
func (c *Options) SetParams(params map[string]string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	c.Fields["PARAMS"] = params
}

// GetParams gets the parameters.
func (c Options) GetParams() map[string]string {
	return c.Fields["PARAMS"]
}

// SetParam sets a parameter.
func (c Options) SetParam(key string, value string) {
	c.SetField("PARAMS", key, value)
}

// GetParam get a parameter.
func (c Options) GetParam(key string) string {
	return c.GetField("PARAMS", key)
}

// SetHeaders re-sets the hole header parameters fields.
func (c *Options) SetHeaders(params map[string]string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	c.Fields["HEADERS"] = params
}

// GetHeaders gets the header parameters.
func (c Options) GetHeaders() map[string]string {
	return c.Fields["HEADERS"]
}

// SetHeader sets a header parameter.
func (c Options) SetHeader(key string, value string) {
	c.SetField("HEADERS", key, value)
}

// GetHeader gets a header parameter.
func (c Options) GetHeader(key string) string {
	return c.GetField("HEADERS", key)
}

// Marshal returns the object to bytes in order to be able to transfer it over HTTP or GRPC or whatever.
func (c Options) Marshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

// Unmarshal converts an Option object from bytes to an actually object.
func (c *Options) Unmarshal(data []byte) *Options {
	if data == nil {
		return NewOptions()
	}
	err := json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	return c
}

// SetField creates a Fields map[string]map[string]string. We need this method in order to initialize the map.
func (c *Options) SetField(group string, key string, value string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	if c.Fields[group] == nil {
		c.Fields[group] = make(map[string]string)
	}
	c.Fields[group][key] = value
}

// GetField returns the value of a Field. The return value is a string.
func (c Options) GetField(group string, key string) string {
	if c.Fields == nil {
		return ""
	}
	if c.Fields[group] == nil {
		return ""
	}
	return c.Fields[group][key]
}
