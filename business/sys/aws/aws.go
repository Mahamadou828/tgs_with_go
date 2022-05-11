//Package aws provide an api to interact with all aws service.
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type AWS struct {
	logger  *zap.SugaredLogger
	sess    *session.Session
	Ssm     *Ssm
	Cognito *Cognito
}

type Config struct {
	CognitoClientID   string
	CognitoUserPoolID string
}

func New(logger *zap.SugaredLogger, cfg Config) (*AWS, error) {
	//Initiate a new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		return nil, err
	}

	return &AWS{
		logger:  logger,
		sess:    sess,
		Ssm:     NewSsm(logger, sess),
		Cognito: NewCognito(logger, sess, cfg.CognitoClientID, cfg.CognitoUserPoolID),
	}, nil
}
