package mocks

// MockSocket is a mock implementation of Socket
// used in the test suite
type MockSocket struct{}

// WriteJSON always returns nil
// used in the test suite
func (s MockSocket) WriteJSON(i interface{}) error {
	return nil
}
