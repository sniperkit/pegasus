package blunder

// Manager describes the error which will handle all the error according to the following specific flags. Some of the
// flag could be importance, priority, store, fatal, etc... The Error manager will send the error at third party
// provider if needed in order to store it or analyze it. Those third party provider could be local micro-services.
type Manager struct {

	// The importance values are specific values which define when and how we will handle the following error.
	importance

	// The error message
	message string

	// The error object if exists
	err error
}

// Set the error message and the error object. The message could pass as string  parameter as the error object. The
// function will return a Manager object. The error manager struct embedded the importance stuck in order to be
// able to set the error. So after e := error.Manager(...) we could use the following practise e.SetFatal() in order
// to set the importance.
func Set(message string, err error) *Manager {
	manager := &Manager{}
	manager.message = message
	manager.err = err
	manager.importance.SetInfo()
	return manager
}

// Handle the error according to struct properties. If the err object is nill the the handle function will silence quit,
// else will check each case of importance.
func (m *Manager) Handle() {

	if m.err == nil {
		return
	}

	switch m.importance.value {
	case 0:
	case 1:
	case 2:
	case 3:
	case 4:
		panic(m.message)
	}
}
