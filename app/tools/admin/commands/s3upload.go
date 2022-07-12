package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
)

//S3UploadJSONFile take a local file and upload it to a s3 bucket. For now we support only json files
func S3UploadJSONFile(cl *aws.Client) error {
	var file, bucket, key string
	fmt.Printf("enter the file path for upload: ")
	if _, err := fmt.Scan(&file); err != nil {
		return fmt.Errorf("invalid secret name: %v", err)
	}
	fmt.Printf("enter the bucket name: ")
	if _, err := fmt.Scan(&bucket); err != nil {
		return fmt.Errorf("invalid secret name: %v", err)
	}
	fmt.Printf("enter the key name: ")
	if _, err := fmt.Scan(&key); err != nil {
		return fmt.Errorf("invalid secret name: %v", err)
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = cl.S3.Upload(b, bucket, key, "application/json")

	return err
}
