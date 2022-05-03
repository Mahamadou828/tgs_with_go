package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func New() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}
