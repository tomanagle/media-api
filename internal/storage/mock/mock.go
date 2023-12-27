package mock

import (
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(session *session.Session, file io.Reader, fileName string) (*s3manager.UploadOutput, error) {
	args := m.Called(session, file, fileName)
	return args.Get(0).(*s3manager.UploadOutput), args.Error(1)
}

func (m *MockStorage) NewSession() (*session.Session, error) {
	args := m.Called()
	return args.Get(0).(*session.Session), args.Error(1)
}
