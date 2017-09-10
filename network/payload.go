package network

type Payload struct {
	Body    []byte
	Options []byte
}

func BuildPayload(body []byte, options []byte) Payload {
	return Payload{Body: body, Options: options}
}
