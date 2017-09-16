package network

import (
	"encoding/json"
)

// Options struct define the parameters that we can pass over network. The path is the main identifier the of the
// options could go at fields property. The Options transferred to the (order side) client always.
type Options struct {

	// The path property could be the RabbitMQ queue ID or a HTTP endpoint.
	Path string

	// Fields are general properties may client need to know
	Fields map[string]map[string]string
}

// NewOptions create a Option object, initialize the struct properties and returns it
func NewOptions() *Options {
	return &Options{Fields: make(map[string]map[string]string)}
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
	err := json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	return c
}

// CreateNewField creates a new property for current Field if not already exists
func (c *Options) CreateNewField(key string) {
	if c.Fields[key] == nil {
		c.Fields[key] = make(map[string]string)
	}
}
