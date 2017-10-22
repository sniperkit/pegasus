package peg

// Channel is the only and main way to transport data from Handlers to receivers. It contains a payload struct field
// with two nested fields payload.Options and payload.Body.
//
// Channel exists only in peg.Handler, peg.Middleware and peg.Client::Listen files as parameter for those
// functions
//
// The payload.Options can be used in order to transport some options via options.Headers or some parameters via
// options.Params. The options.Body is similar with http body and is used in order to transfer raw content.
//
// An example for usage could be:
//	func handler(channel *peg.Channel) {
//		payload := channel.Receive() // Receive the payload here
//		options := peg.BuildOptions(payload.Options) // Convert the received options to options struct
// 		payload.Body = payload.Body + []byte(" sub-fix") // Change the body
//	}
type Channel struct {
	payload chan Payload
}

// NewChannel creates a new Channel which contains a new payload object.
//
// The channel size has to be set as parameter.
//
// The channel size defines how big the channel will be. The buffer size is the number of elements that can be sent to
// the channel without the send blocking. By default, a channel has a buffer size of 0
// (you get this with make(chan int)). This means that every single send will block until another goroutine receives
// from the channel.
func NewChannel(channelSize int) *Channel {
	return &Channel{payload: make(chan Payload, channelSize)}
}

// Send method get a payload object and send it through the channel.
func (c Channel) Send(payload Payload) {
	c.payload <- payload
}

// Receive data from channel and return it. The function return a Payload.
func (c Channel) Receive() Payload {
	return <-c.payload
}
