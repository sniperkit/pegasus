package network

// Handler is a type of function which has as parameter a channel object. Usually it uses for transporters.
type Handler func(chanel *Channel)
