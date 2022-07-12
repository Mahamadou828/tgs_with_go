package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SSM provide an api to interact with the
//aws simple secret manager
type SSM struct {
	svc     *secretsmanager.SecretsManager
	env     string
	service string
}

type Secret struct {
	Name  string
	Value string
}

func NewSSM(sess *session.Session, service, env string) *SSM {
	svc := secretsmanager.New(sess)
	return &SSM{
		svc:     svc,
		service: service,
		env:     env,
	}
}

//ListSecrets Retrieve all secrets store in the aws account
//and filter them based on the service pass and the build.
func (s *SSM) ListSecrets() (map[string]string, error) {
	var secrets map[string]string

	result, err := s.svc.ListSecrets(
		&secretsmanager.ListSecretsInput{
			Filters: []*secretsmanager.Filter{
				{
					Key:    aws.String(secretsmanager.FilterNameStringTypeTagKey),
					Values: []*string{aws.String("service"), aws.String("env")},
				},
				{
					Key:    aws.String(secretsmanager.FilterNameStringTypeTagValue),
					Values: []*string{aws.String(s.service), aws.String(s.env)},
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("can't get secrets for service %s in environment %s: %v", s.service, s.env, err)
	}

	if len(result.SecretList) == 0 {
		return nil, fmt.Errorf("no pool found for service %s in environment %s", s.service, s.env)
	}

	secretVal, err := s.svc.GetSecretValue(
		&secretsmanager.GetSecretValueInput{
			SecretId: result.SecretList[0].Name,
		},
	)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *secretVal.Name, secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *secretVal.Name, secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *secretVal.Name, secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *secretVal.Name, secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				return nil, fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *secretVal.Name, secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				return nil, fmt.Errorf(aerr.Error())
			}
		} else {
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(*secretVal.SecretString), &secrets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %s, error: %v", *secretVal.Name, err)
	}

	return secrets, nil
}

//GetPoolSecrets Retrieve all secrets store for a specified pool name.
//No tag check or env check would be performed on the pool.
func (s *SSM) GetPoolSecrets(poolName string, dest any) error {
	result, err := s.svc.GetSecretValue(
		&secretsmanager.GetSecretValueInput{SecretId: aws.String(poolName)},
	)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *result.Name, secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				return fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *result.Name, secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *result.Name, secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				return fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *result.Name, secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				return fmt.Errorf("failed to retrieve secret: %s, error: %s, %s", *result.Name, secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			return err
		}
	}

	if err := json.Unmarshal([]byte(*result.SecretString), dest); err != nil {
		return fmt.Errorf("failed to unmarshal secret: %s, error: %v", *result.Name, err)
	}

	return nil
}

//CreatePool create a new secret pool, a secret pool is a set of secret for
//a given service. If a pool already exists, no change is made and no error is returned.
func (s *SSM) CreatePool() error {
	m := make(map[string]string)
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("can't create secret pool: %v", err)
	}
	input := &secretsmanager.CreateSecretInput{
		Name:        aws.String(fmt.Sprintf("%s/%s", s.service, s.env)),
		Description: aws.String(fmt.Sprintf("secret pool for service %s in %s", s.service, s.env)),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String("service"),
				Value: aws.String(s.service),
			},
			{
				Key:   aws.String("env"),
				Value: aws.String(s.env),
			},
		},
		SecretString: aws.String(string(b)),
	}

	_, err = s.svc.CreateSecret(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceExistsException:
				return nil
			case secretsmanager.ErrCodeInvalidParameterException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				return fmt.Errorf("error creating secret: %s, %s", secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
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

//CreateSecrets insert new secret in the secret pool. If the given secret already exists
//its value will be updated.
func (s *SSM) CreateSecrets(sts []Secret) error {
	secrets, err := s.ListSecrets()
	if err != nil {
		return err
	}
	for _, st := range sts {
		name, value := st.Name, st.Value
		//placing the secret inside the map
		//if the secret is already created, this will overwrite it's value
		secrets[name] = value
	}

	b, err := json.Marshal(secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal the secret: %v", err)
	}

	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(fmt.Sprintf("%s/%s", s.service, s.env)),
		SecretString: aws.String(string(b)),
	}

	_, err = s.svc.UpdateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeLimitExceededException:
				return errors.New(secretsmanager.ErrCodeLimitExceededException)
			case secretsmanager.ErrCodeEncryptionFailure:
				return errors.New(secretsmanager.ErrCodeEncryptionFailure)
			case secretsmanager.ErrCodeResourceExistsException:
				return errors.New(secretsmanager.ErrCodeResourceExistsException)
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				return errors.New(secretsmanager.ErrCodeMalformedPolicyDocumentException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			case secretsmanager.ErrCodePreconditionNotMetException:
				return errors.New(secretsmanager.ErrCodePreconditionNotMetException)
			case secretsmanager.ErrCodeDecryptionFailure:
				return errors.New(secretsmanager.ErrCodeDecryptionFailure)
			default:
				return errors.New(aerr.Error())
			}
		}
	}
	return nil

}

//CreateSecret insert a new secret in the secret pool. If the given secret already exists
//its value will be updated.
func (s *SSM) CreateSecret(name, value string) error {
	secrets, err := s.ListSecrets()
	if err != nil {
		return err
	}

	//placing the secret inside the map
	//if the secret is already created, this will overwrite it's value
	secrets[name] = value

	b, err := json.Marshal(secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal the secret: %v", err)
	}

	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(fmt.Sprintf("%s/%s", s.service, s.env)),
		SecretString: aws.String(string(b)),
	}

	_, err = s.svc.UpdateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeLimitExceededException:
				return errors.New(secretsmanager.ErrCodeLimitExceededException)
			case secretsmanager.ErrCodeEncryptionFailure:
				return errors.New(secretsmanager.ErrCodeEncryptionFailure)
			case secretsmanager.ErrCodeResourceExistsException:
				return errors.New(secretsmanager.ErrCodeResourceExistsException)
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				return errors.New(secretsmanager.ErrCodeMalformedPolicyDocumentException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			case secretsmanager.ErrCodePreconditionNotMetException:
				return errors.New(secretsmanager.ErrCodePreconditionNotMetException)
			case secretsmanager.ErrCodeDecryptionFailure:
				return errors.New(secretsmanager.ErrCodeDecryptionFailure)
			default:
				return errors.New(aerr.Error())
			}
		}
	}

	return nil
}
