package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type S3 struct {
	DownloaderManager *s3manager.Downloader
	UploaderManager   *s3manager.Uploader
	Log               *zap.SugaredLogger
}

func NewS3(log *zap.SugaredLogger, sess *session.Session) *S3 {
	return &S3{
		DownloaderManager: s3manager.NewDownloader(sess),
		UploaderManager:   s3manager.NewUploader(sess),
		Log:               log,
	}
}

func (s S3) Download(w io.WriterAt, bucketName, key string) (int64, error) {
	numBytes, err := s.DownloaderManager.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, err
	}

	return numBytes, nil
}
