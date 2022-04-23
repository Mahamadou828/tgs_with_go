//Package ssm provide an interface to access the aws simple secret manager service
//For more details, see: https://docs.aws.amazon.com/sdk-for-go/api/
package ssm

//@todo pass the session in params

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//ListSecrets Retrieve all secrets bound to that specific account
//and filter them based on the service pass and the build.
func ListSecrets(service string, build string) (map[string]string, error) {
	//Initiate a new aws session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		return nil, err
	}

	svc := secretsmanager.New(sess)

	input := &secretsmanager.ListSecretsInput{
		Filters: []*secretsmanager.Filter{
			{
				Key:    aws.String(secretsmanager.FilterNameStringTypeTagKey),
				Values: []*string{aws.String("service"), aws.String("build")},
			},
			{
				Key:    aws.String(secretsmanager.FilterNameStringTypeTagValue),
				Values: []*string{aws.String(service), aws.String(build)},
			},
		},
	}

	result, err := svc.ListSecrets(input)

	secrets := make(map[string]string)

	for _, value := range result.SecretList {
		input := &secretsmanager.GetSecretValueInput{
			SecretId: value.Name,
		}

		result, err := svc.GetSecretValue(input)

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

//CreateSecret creates a new secret and host it inside the aws ssm service
func CreateSecret(name string, value string, service string, build string, desc string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	})

	if err != nil {
		//@todo handle the error
		panic(err)
	}

	svc := secretsmanager.New(sess)

	input := &secretsmanager.CreateSecretInput{
		Description: aws.String(desc),
		Name:        aws.String(name),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String("service"),
				Value: aws.String(service),
			},
			{
				Key:   aws.String("build"),
				Value: aws.String(build),
			},
		},
		SecretString: aws.String(value),
	}

	_, err = svc.CreateSecret(input)

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
