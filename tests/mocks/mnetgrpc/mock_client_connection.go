package mnetgrpc

// Close mock for grpc.ClientConn struct
type MockClientConnection struct {
	CloseMock func() error
}

// Close mock for grpc.Close function
func (m MockClientConnection) Close() error {
	if m.CloseMock != nil {
		return m.CloseMock()
	}
	return nil
}
