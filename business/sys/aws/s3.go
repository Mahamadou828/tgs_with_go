package aws

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsToolkit "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type S3 struct {
	downloaderManager *s3manager.Downloader
	uploaderManager   *s3manager.Uploader
	Log               *zap.SugaredLogger
}

func NewS3(log *zap.SugaredLogger, sess *session.Session) *S3 {
	return &S3{
		downloaderManager: s3manager.NewDownloader(sess),
		uploaderManager:   s3manager.NewUploader(sess),
		Log:               log,
	}
}

func (s S3) Download(w io.WriterAt, bucketName, key string) (int64, error) {
	numBytes, err := s.downloaderManager.Download(w, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, err
	}

	return numBytes, nil
}

func (s S3) Upload(b []byte, bucketName, key, env, contentType string) error {
	_, err := s.uploaderManager.Upload(&s3manager.UploadInput{
		Body:        strings.NewReader(string(b)),
		Bucket:      awsToolkit.String(bucketName),
		ContentType: awsToolkit.String(contentType),
		Key:         awsToolkit.String(key),
		Metadata:    map[string]*string{"env": awsToolkit.String(env)},
	})

	if err != nil {
		return err
	}

	return nil
}
