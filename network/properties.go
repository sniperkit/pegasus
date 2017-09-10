package network

type Properties struct {
	Path    string
	Method string
	Fields map[string]map[string]string
	Objects map[string]map[string]interface{}
}

func NewProperties(path string, method string) *Properties {
	return &Properties{Path: path, Method: method}
}