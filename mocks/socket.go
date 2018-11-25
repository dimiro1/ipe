package mocks

// MockSocket is a mock implementation of Socket
// used in the test suite
type MockSocket struct{}

func (s MockSocket) WriteJSON(i interface{}) error {
	return nil
}
