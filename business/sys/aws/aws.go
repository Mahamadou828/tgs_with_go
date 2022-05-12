//Package aws provide an api to interact with all aws service.
package aws

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
	"os"
)

type AWS struct {
	logger  *zap.SugaredLogger
	sess    *session.Session
	Ssm     *Ssm
	Cognito *Cognito
}

type Config struct {
	Account string
	Service string
	Env     string
}

type parser struct {
	Secrets map[string]string
}

func New(logger *zap.SugaredLogger, config Config) (*AWS, error) {

	_, err := os.OpenFile("/root/.aws/credentials", os.O_RDWR, 0755)

	//Initiate a new aws session
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                        aws.String("eu-west-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		Profile: config.Account,
	})

	if err != nil {
		return nil, err
	}

	ssm := NewSsm(logger, sess)

	cfg := struct {
		Cognito struct {
			UserPoolID string `conf:"required"`
			ClientID   string `conf:"required"`
		}
	}{}

	if err := extractConfigFromSecrets(ssm, config.Service, config.Env, &cfg); err != nil {
		return nil, err
	}

	return &AWS{
		logger:  logger,
		sess:    sess,
		Ssm:     NewSsm(logger, sess),
		Cognito: NewCognito(logger, sess, cfg.Cognito.ClientID, cfg.Cognito.UserPoolID),
	}, nil
}

//extractConfigFromSecrets Use ssm to extract all config needed to
//start other aws services ( cognito, s3 etc )
func extractConfigFromSecrets(ssm *Ssm, service, env string, cfg any) error {
	secrets, err := ssm.ListSecrets(service, env)

	if err != nil {
		return err
	}

	if _, err := config.Parse(cfg, service, parser{Secrets: secrets}); err != nil {
		return err
	}

	return nil
}

func (p parser) Parse(field config.Field) error {
	val, ok := p.Secrets[field.Name]

	if !ok {
		return fmt.Errorf("can't find configuration for aws service. Missing secrets for: %s", field.Name)
	}

	return config.SetFieldValue(field, val)
}
