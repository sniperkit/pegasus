package grpc_transport

import (
	"bitbucket.org/code_horse/pegasus/network"
	pb "bitbucket.org/code_horse/pegasus/transport/grpc_transport/proto"
)

type Properties struct {
	network.Properties
}

func NewProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	properties.Objects = make(map[string]map[string]interface{})
	return properties
}

func (p *Properties) BuildProperties(properties *network.Properties) *Properties {
	p.Properties = *properties
	return p
}

func (p *Properties) GetProperties() *network.Properties {
	return &p.Properties
}

func (p *Properties) initMapperFields(key string) {
	if p.Fields[key] == nil {
		p.Fields[key] = make(map[string]string)
	}
}

func (p *Properties) initMapperObjects(key string) {
	if p.Fields[key] == nil {
		p.Objects[key] = make(map[string]interface{})
	}
}

func (p *Properties) SetPath(path string) *Properties {
	p.Path = path
	return p
}

func (p *Properties) GetPath() string {
	return p.Path
}

func (p *Properties) SetConnection(connection pb.ServeClient) *Properties {
	p.initMapperObjects("GRPC-CONNECTION")
	p.Objects["GRPC-CONNECTION"]["VALUE"] = connection
	return p
}

func (p *Properties) GetConnection() pb.ServeClient{
	if c, ok := p.Objects["GRPC-CONNECTION"]["VALUE"].(pb.ServeClient); ok {
		return c
	}
	return nil
}