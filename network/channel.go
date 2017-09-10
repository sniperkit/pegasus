package network

type Channel struct {
	payload chan Payload
}

func NewChannel(channelSize int) *Channel {
	return &Channel{payload:make(chan Payload, channelSize)}
}

func (c *Channel) Send(payload Payload) {
	c.payload <- payload
}

func (c *Channel) Receive() Payload {
	return <- c.payload
}
