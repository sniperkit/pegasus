package network

// Properties struct contains all the information which don't need to send in the order side (client). Some of those
// could be the RabbitMQ Qos or Queue declaration or the path and the method for HTTP request handlers
type Properties struct {

	// Path is the identifier for handler. For example the /products/all for HTTP handler
	Path string

	// todo: [fix] [code:A001] check if we could to remove it
	// The method is really used for HTTP handlers
	Method string

	// Fields are data and properties that we can set
	Fields map[string]map[string]string

	// Object usually are hole struct object or whatever else we wants.
	Objects map[string]map[string]interface{}
}

// NewProperties generate and return a Property object
// todo: [fix] [code:A001] check if we could remove the method parameter
func NewProperties(path string, method string) *Properties {
	return &Properties{
		Path: path,
		Method: method,
		Fields: make(map[string]map[string]string),
		Objects: make(map[string]map[string]interface{}),
	}
}

// CreateNewField creates a new property for current Fields if not already exists
func (p *Properties) CreateNewField(key string) {
	if p.Fields[key] == nil {
		p.Fields[key] = make(map[string]string)
	}
}

// CreateNewObject create a new property for current Objects if not already exists
func (p *Properties) CreateNewObject(key string) {
	if p.Objects[key] == nil {
		p.Objects[key] = make(map[string]interface{})
	}
}
