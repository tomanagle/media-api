package storage

import (
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage interface {
	Upload(session *session.Session, file io.Reader, fileName string) (*s3manager.UploadOutput, error)
	NewSession() (*session.Session, error)
}
