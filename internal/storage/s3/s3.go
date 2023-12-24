package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	bucketName string
	awsRegion  string
}

type NewS3Params struct {
	BucketName string
	AWSRegion  string
}

func New(p NewS3Params) *S3 {
	return &S3{
		bucketName: p.BucketName,
		awsRegion:  p.AWSRegion,
	}
}

func (s *S3) Upload(session *session.Session, file io.Reader, fileName string) (*s3manager.UploadOutput, error) {

	// S3 service client the Upload manager will use.
	s3Svc := s3.New(session)

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	upParams := &s3manager.UploadInput{
		Bucket: &s.bucketName,
		Key:    &fileName,
		Body:   file,
	}

	result, err := uploader.Upload(upParams)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *S3) NewSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(s.awsRegion)},
	)
}
