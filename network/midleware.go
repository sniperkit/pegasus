package network

type Middleware func(handler Handler, chanel *Channel)
