package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.uber.org/zap"
)

// Ssm provide an api to interact with the
//aws simple secret manager
type Ssm struct {
	svc    *secretsmanager.SecretsManager
	logger *zap.SugaredLogger
}

func NewSsm(logger *zap.SugaredLogger, sess *session.Session) *Ssm {
	svc := secretsmanager.New(sess)
	return &Ssm{
		svc:    svc,
		logger: logger,
	}
}

//ListSecrets Retrieve all secrets store in the aws account
//and filter them based on the service pass and the build.
func (s *Ssm) ListSecrets(service, env string) (map[string]string, error) {
	input := &secretsmanager.ListSecretsInput{
		Filters: []*secretsmanager.Filter{
			{
				Key:    aws.String(secretsmanager.FilterNameStringTypeTagKey),
				Values: []*string{aws.String("service"), aws.String("env")},
			},
			{
				Key:    aws.String(secretsmanager.FilterNameStringTypeTagValue),
				Values: []*string{aws.String(service), aws.String(env)},
			},
		},
	}

	result, err := s.svc.ListSecrets(input)

	if err != nil {
		return nil, err
	}

	secrets := make(map[string]string)

	for _, value := range result.SecretList {
		input := &secretsmanager.GetSecretValueInput{
			SecretId: value.Name,
		}

		result, err := s.svc.GetSecretValue(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case secretsmanager.ErrCodeResourceNotFoundException:
					return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *value.Name, secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				case secretsmanager.ErrCodeInvalidParameterException:
					return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *value.Name, secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				case secretsmanager.ErrCodeInvalidRequestException:
					return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *value.Name, secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				case secretsmanager.ErrCodeDecryptionFailure:
					return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *value.Name, secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				case secretsmanager.ErrCodeInternalServiceError:
					return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *value.Name, secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				default:
					return nil, fmt.Errorf(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				return nil, err
			}
		}

		secrets[*result.Name] = *result.SecretString
	}

	return secrets, nil
}

func (s *Ssm) CreateSecret(name, value, service, env, desc string) error {
	input := &secretsmanager.CreateSecretInput{
		Description: aws.String(desc),
		Name:        aws.String(name),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String("service"),
				Value: aws.String(service),
			},
			{
				Key:   aws.String("env"),
				Value: aws.String(env),
			},
		},
		SecretString: aws.String(value),
	}

	_, err := s.svc.CreateSecret(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			case secretsmanager.ErrCodePreconditionNotMetException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			default:
				return fmt.Errorf("error creating secret: %s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return err
		}
	}

	return nil
}
