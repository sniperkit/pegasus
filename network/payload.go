package network

type Payload struct {
	Body    []byte
	Options []byte
}

func NewPayload(body []byte, options []byte) *Payload {
	return &Payload{Body: body, Options: options}
}

func BuildPayload(body []byte, options []byte) Payload {
	return Payload{Body: body, Options: options}
}
