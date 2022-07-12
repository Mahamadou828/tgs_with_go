//Package aws provide an api to interact with all aws service.
package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

//Client provides an api to interact with AWS services
type Client struct {
	logger  *zap.SugaredLogger
	sess    *session.Session
	SSM     *SSM
	Cognito *Cognito
	S3      *S3
	env     string
	service string
	account string
}

type Config struct {
	Service string
	Env     string
	Cognito struct {
		UserPoolID string
		ClientID   string
		//Seed is used to generate unique sub that are used as username
		//for cognito user
		Seed string
	}
}

func New(logger *zap.SugaredLogger, cfg Config) (*Client, error) {

	//Initiate a new aws session
	sess, err := createSess()
	if err != nil {
		return nil, fmt.Errorf("failed to initiate a new session: %v", err)
	}

	return &Client{
		logger:  logger,
		sess:    sess,
		SSM:     NewSSM(sess, cfg.Service, cfg.Env),
		Cognito: NewCognito(logger, sess, cfg.Cognito.ClientID, cfg.Cognito.UserPoolID, cfg.Cognito.Seed),
		S3:      NewS3(logger, sess, cfg.Service, cfg.Env),
		env:     cfg.Env,
		service: cfg.Service,
	}, nil
}

//GetSecretList allow to retrieve the secret for a given service in a given environment
//without having to create a new aws client. This function is aimed to be used only for
//configuration purposes.
func GetSecretList(service, env string) (map[string]string, error) {
	//Initiate a new aws session
	sess, err := createSess()
	if err != nil {
		return nil, fmt.Errorf("failed to initiate a new session: %v", err)
	}

	return NewSSM(sess, service, env).ListSecrets()
}

func createSess() (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                        aws.String("eu-west-1"),
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
		Profile: "formation",
	})
}