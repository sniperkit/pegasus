package network

import (
	"encoding/json"
)

// Options struct define the parameters that we can pass over network. The path is the main identifier the of the
// options could go at fields property. The Options transferred to the (order side) client always.
type Options struct {

	// Fields are general properties may client need to know
	Fields map[string]map[string]string
}

// NewOptions create a Option object, initialize the struct properties and returns it
func NewOptions() *Options {
	return &Options{Fields: make(map[string]map[string]string)}
}

// BuildOptions receives bytes as parameter and convert them to Option object.
func BuildOptions(data []byte) *Options {
	return NewOptions().Unmarshal(data)
}

// SetParams re-set the hole parameters fields
func (c *Options) SetParams(params map[string]string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	c.Fields["PARAMS"] = params
}

// GetParams get the parameters
func (c *Options) GetParams() map[string]string {
	return c.Fields["PARAMS"]
}

// SetParam set a parameter
func (c *Options) SetParam(key string, value string) {
	c.SetField("PARAMS", key, value)
}

// GetParam get a parameter
func (c *Options) GetParam(key string) string {
	return c.GetField("PARAMS", key)
}

// SetHeaders re-set the hole header parameters fields
func (c *Options) SetHeaders(params map[string]string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	c.Fields["HEADERS"] = params
}

// GetHeaders get the header parameters
func (c *Options) GetHeaders() map[string]string {
	return c.Fields["HEADERS"]
}

// SetHeader set a header parameter
func (c *Options) SetHeader(key string, value string) {
	c.SetField("HEADERS", key, value)
}

// GetHeader get a header parameter
func (c *Options) GetHeader(key string) string {
	return c.GetField("HEADERS", key)
}

// Marshal return the hole object to byte in order to be able to transfer it over HTTP or GRPC or whatever
func (c *Options) Marshal() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

// Unmarshal convert an Option object from bytes to an actually object.
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

// SetField create a Fields map[string]map[string]string. We need this method in order to initialize the map each
// time when a new group (the first map is known as group) is created. The parameters are the group which is a key
// for the first mapper, the key which for the previous group and the value.
func (c *Options) SetField(group string, key string, value string) {
	if c.Fields == nil {
		c.Fields = make(map[string]map[string]string)
	}
	if c.Fields[group] == nil {
		c.Fields[group] = make(map[string]string)
	}
	c.Fields[group][key] = value
}

// GetField return the value of a Field. The return value is a string.
func (c *Options) GetField(group string, key string) string {
	if c.Fields == nil {
		return ""
	}
	if c.Fields[group] == nil {
		return ""
	}
	return c.Fields[group][key]
}
