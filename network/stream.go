package network

type IStream interface {
	Send(payload *Payload)
	Receive() (*Payload, error)
}
