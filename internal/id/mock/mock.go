package mock

import "github.com/stretchr/testify/mock"

type MockID struct {
	mock.Mock
}

func (m *MockID) New() (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}
