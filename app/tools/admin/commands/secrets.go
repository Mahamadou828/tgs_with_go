package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	awsToolkit "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

const (
	S3Src        = "s3file"
	CliSrc       = "cli"
	LocalFileSrc = "localfile"
)

type CreateSecretCfg struct {
	AwsConfig aws.Config
	SrcType   string
	Log       *zap.SugaredLogger
	Filename  string
	Service   string
	Env       string
	Bucket    string
	Key       string
	Secrets   []string
}

type Secret struct {
	Name        string `json:"Name"`
	Value       string `json:"Value"`
	Description string `json:"Description"`
}

//CreateSecret create new secret inside aws secret manager service.
//It's possible to specify three secret src: s3 bucket file, local
//file and cli param. The format of file should be json.
//And the format of cli param should be --secrets=["secretname:secretvalue:secretdesc"]
func CreateSecret(cfg CreateSecretCfg) error {
	var secrets []Secret
	sessAws, err := aws.New(cfg.Log, cfg.AwsConfig)

	if err != nil {
		return err
	}

	switch cfg.SrcType {
	case S3Src:
		buffer := awsToolkit.NewWriteAtBuffer([]byte{})
		_, err := sessAws.S3.DownloaderManager.Download(buffer, &s3.GetObjectInput{
			Bucket: awsToolkit.String("tgs-with-go-secrets"),
			Key:    awsToolkit.String("secrets-development-tgs-api"),
		})
		if err != nil {
			return err
		}
		buf := buffer.Bytes()
		if len(buf) == 0 {
			return fmt.Errorf("downloaded file is empty")
		}
		fmt.Println(string(buf))
		if err := json.Unmarshal(buf, &secrets); err != nil {
			return err
		}
	case CliSrc:
		for _, secret := range cfg.Secrets {
			parseSecret := strings.Split(secret, ":")
			if len(parseSecret) != 3 {
				return fmt.Errorf("secrets malformated: %v, secret should be in the following format: key:value:description", secret)
			}
			name, value, desc := parseSecret[0], parseSecret[1], parseSecret[2]
			secrets = append(secrets, Secret{
				Name:        name,
				Value:       value,
				Description: desc,
			})
		}
	case LocalFileSrc:
		jsonFile, err := os.Open(cfg.Filename)
		if err != nil {
			return err
		}
		buf, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(buf, &secrets); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown source type %s", cfg.SrcType)
	}

	for _, secret := range secrets {
		if err := sessAws.Ssm.CreateSecret(secret.Name, secret.Value, cfg.Service, cfg.Env, secret.Description); err != nil {
			return err
		}
	}

	return nil
}
