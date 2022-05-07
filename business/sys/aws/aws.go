//Package aws provide an api to interact with all aws service.
package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type AWS struct {
	logger *zap.SugaredLogger
	sess   *session.Session
	Ssm    *Ssm
}

func New(logger *zap.SugaredLogger) (*AWS, error) {
	//Initiate a new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		return nil, err
	}

	return &AWS{
		logger: logger,
		sess:   sess,
		Ssm:    NewSsm(logger, sess),
	}, nil
}
