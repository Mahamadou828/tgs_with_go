package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
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
