package network

// Channel is the only and main way to transport data between transporters. The payload property contains the body and
// options as categories for transportation data.
type Channel struct {
	payload chan Payload
}

// NewChannel create a new Channel which contains a new payload object. The channel size have to be set as parameter.
func NewChannel(channelSize int) *Channel {
	return &Channel{payload: make(chan Payload, channelSize)}
}

// Send method get a payload object and send it through the channel.
func (c *Channel) Send(payload Payload) {
	c.payload <- payload
}

// Receive data from channel and return it.
func (c *Channel) Receive() Payload {
	return <-c.payload
}
