package blunder

// Error importance describes how important an error could be for us and what actions should we do for that.
type importance struct {

	// The value could get the the number between 0-4
	//
	// value: 0 is the default value. The system will store the error.
	//
	// value: 1 is the warning. Something is wrong and a problem may occur. System stores this error.
	//
	// value: 2 is the error. Something is broken like an external HTTP call or param validation. System stores the
	// 			error message and stops the execution of the current call.
	//
	// value: 3 is the high importance error. Something goes wrong and is related with sensitive data. Could be a
	// 			Payment transaction or wrong configuration for pricing, etc.. System stores this error and stops the
	// 			execution of the current call.
	//
	// value: 4 is the panic error. Something happens and the application is unable to start.We stop it and we
	//           throw a panic alert.
	value int
}

// SetInfo sets value: 0 is the default value. The system will store the error.
func (i *importance) SetInfo() {
	i.value = 0
}

// SetWarning sets value: 1 is the warning. Something is wrong and a problem may occur. System stores this error.
func (i *importance) SetWarning() {
	i.value = 1
}

// SetError sets value: 2 is the error. Something is broken like an external HTTP call or param validation. System stores the
// error message and stops the execution of the current call.
func (i *importance) SetError() {
	i.value = 2
}

// SetFatal sets value: 3 is the high importance error. Something goes wrong and is related with sensitive data. Could be a
// Payment transaction or wrong configuration for pricing, etc.. System stores this error and stops the
// execution of the current call.
func (i *importance) SetFatal() {
	i.value = 3
}

// SetPanicError sets value: 4 is the panic error. Something happens and the application is unable to start.We stop it and we
// throw a panic alert.
func (i *importance) SetPanicError() {
	i.value = 4
}
