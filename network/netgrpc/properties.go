package netgrpc

import (
	"bitbucket.org/code_horse/pegasus/network"
	"google.golang.org/grpc"
)

// Properties container embedded struct network.Properties. It's an exaction of properties designed for GRPC protocol.
type Properties struct {
	network.Properties
}

// NewProperties initialize the fields and return a new Properties object.
func NewProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	properties.Objects = make(map[string]map[string]interface{})
	return properties
}

// SetPath sets the handler unique path
func (p *Properties) SetPath(path string) *Properties {
	p.Path = path
	return p
}

// GetPath return the handler unique path
func (p *Properties) GetPath() string {
	return p.Path
}
