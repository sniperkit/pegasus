package network

// Channel is the only and main way to transport data between from Handlers to receivers. The payload struct field
// contains the Body and Options as sub-struct fields.
//
// An example for usage could be:
//	func handler(channel *network.Channel) {
//		payload := channel.Receive() // Receive the payload here
//		options := network.BuildOptions(payload.Options) // Convert the received options to options struct
// 		payload.Body = payload.Body + []byte(" sub-fix") // Change the body
//
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
