package http_transport

import "bitbucket.org/code_horse/pegasus/network"

type Properties struct {
	network.Properties
}

func NewProperties() *Properties {
	properties := &Properties{}
	properties.Fields = make(map[string]map[string]string)
	return properties
}

func (p *Properties) BuildProperties(properties *network.Properties) *Properties {
	return &Properties{Properties: *properties}
}

func (p *Properties) GetProperties() *network.Properties {
	return &p.Properties
}

func (p *Properties) SetPath(path string) *Properties {
	p.Path = path
	return p
}

func (p *Properties) GetPath() string {
	return p.Path
}

func (p *Properties) SetGetMethod() *Properties {
	if p.Fields["METHOD"] == nil {
		p.Fields["METHOD"] = make(map[string]string)
	}
	p.Fields["METHOD"]["VALUE"] = "GET"
	return p
}

func (p *Properties) SetPostMethod() *Properties {
	if p.Fields["METHOD"] == nil {
		p.Fields["METHOD"] = make(map[string]string)
	}
	p.Fields["METHOD"]["VALUE"] = "POST"
	return p
}

func (p *Properties) SetPutMethod() *Properties {
	if p.Fields["METHOD"] == nil {
		p.Fields["METHOD"] = make(map[string]string)
	}
	p.Fields["METHOD"]["VALUE"] = "PUT"
	return p
}

func (p *Properties) SetPatchMethod() *Properties {
	if p.Fields["METHOD"] == nil {
		p.Fields["METHOD"] = make(map[string]string)
	}
	p.Fields["METHOD"]["VALUE"] = "PATCH"
	return p
}

func (p *Properties) SetDeleteMethod() *Properties {
	if p.Fields["METHOD"] == nil {
		p.Fields["METHOD"] = make(map[string]string)
	}
	p.Fields["METHOD"]["VALUE"] = "DELETE"
	return p
}

func (p *Properties) GetMethod() string {
	return p.Fields["METHOD"]["VALUE"]
}

