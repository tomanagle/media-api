package mock

import (
	"go-media/internal/store"

	"github.com/stretchr/testify/mock"
)

type MockMediaStore struct {
	mock.Mock
}

func (m *MockMediaStore) CreateMedia(params store.CreateMediaParams) (media *store.Media, err error) {
	args := m.Called(params)
	return args.Get(0).(*store.Media), args.Error(1)
}
