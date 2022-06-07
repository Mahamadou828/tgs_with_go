package commands

import (
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	awsToolkit "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
	"os"
)

//Upload take a local file and upload it to a s3 bucket. For now we support only json files
func Upload(cfg aws.Config, log *zap.SugaredLogger, file, bucket, key string) error {
	sessAws, err := aws.New(log, cfg)
	if err != nil {
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = sessAws.S3.UploaderManager.Upload(&s3manager.UploadInput{
		Body:        f,
		Bucket:      awsToolkit.String(bucket),
		ContentType: awsToolkit.String("application/json"),
		Key:         awsToolkit.String(key),
		Metadata:    map[string]*string{"env": awsToolkit.String(cfg.Env)},
	})

	return err
}
