package mock

import (
	"context"
	"go-media/internal/store"

	"github.com/stretchr/testify/mock"
)

type MockMediaStore struct {
	mock.Mock
}

func (m *MockMediaStore) CreateMedia(ctx context.Context, params store.CreateMediaParams) (media *store.Media, err error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*store.Media), args.Error(1)
}

func (m *MockMediaStore) GetMedia(ctx context.Context, params store.GetMediaParams) (media []store.Media, err error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]store.Media), args.Error(1)
}
