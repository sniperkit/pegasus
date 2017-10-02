package blunder

// Manager struct describes the errors, the actions after an error occurs and the tracking of it. Manager struct
// contains the following fields
//
// importance: Defines how important an error could be and how to handle it. It's an embedded struct that helps to
// prioritize the error.
//
// message: Defines the error message
//
// error: Is the actual error
//
// Example:
//  err := blunder.Set("This is the message", errors.New("sample error"))
//  err.SetPanicError()
type Manager struct {

	// The importance values are specific values which define when and how we will handle the following error.
	importance

	// The error message
	message string

	// The error object if exists
	err error
}

// Set the error message and the error object. It gets two parameters: the message parameter as string and the error
// object. The function will return a Manager object.
//
// The Manager struct embeds the importance struct in order to be able to prioritize it.
func Set(message string, err error) *Manager {
	manager := &Manager{}
	manager.message = message
	manager.err = err
	manager.importance.SetInfo()
	return manager
}

// Handle the error according to properties. If the Manger.err object is nil the handle function will silently quit,
// otherwise it will check each case of importance stat
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
