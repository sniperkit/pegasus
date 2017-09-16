package blunder

// Error importance describes how important an error could be for us and what actions should we do for that.
type importance struct {

	// The value could get the the number between 0-4
	//
	// value: 0 is the default values and is mean that we only will send the error in order to store it and nothing
	// 			else will happen.
	//
	// value: 1 is the warring. Something is wrong and may accrue a problem. We send this error and store it.
	//
	// value: 2 is the error. Something is broken could be a http call with other service, could be a validation param,
	// 			etc. We store this message and stop the execute of the current call.
	//
	// value: 3 is the high important error. Something is
	//			go wrong and it has to be with very sensitive data. Could be a Payment transaction or wrong
	// 			configurations for pricing, etc.. We store this error and stop the execution of current call.
	// 			Also we notify the programmers immediately.
	//
	// value: 4 is the panic error. Something happens and the application is unable to start. We stop it and then we
	// 			throw a panic alert.
	value int
}

// Is the info we only will send the error in order to store it and nothing else will happen.
func (i *importance) SetInfo() {
	i.value = 0
}

// Something is wrong and may accrue a problem. We send this error and store it.
func (i *importance) SetWarring() {
	i.value = 1
}

// Something is broken could be a http call with other service, could be a validation param, etc. We store this message
// and stop the execute of the current call.
func (i *importance) SetError() {
	i.value = 2
}

// Something is	go wrong and it has to be with very sensitive data. Could be a Payment transaction or wrong
// configurations for pricing, etc.. We store this error and stop the execution of current call. Also we notify
// the programmers immediately.
func (i *importance) SetFatal() {
	i.value = 3
}

// Something happens and the application is unable to start. We stop it and then we throw a panic alert.
func (i *importance) SetPanicError() {
	i.value = 4
}
