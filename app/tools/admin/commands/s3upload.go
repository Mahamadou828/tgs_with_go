package commands

import (
	"io/ioutil"
	"os"

	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"go.uber.org/zap"
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
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = sessAws.S3.Upload(b, bucket, key, cfg.Env, "application/json")

	return err
}
