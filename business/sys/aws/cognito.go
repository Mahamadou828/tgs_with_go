package aws

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//Cognito represent an instance of the cognito session
type Cognito struct {
	log              *zap.SugaredLogger
	identityProvider *cognitoidentityprovider.CognitoIdentityProvider
	clientID         string
	userPoolID       string
}

type Session struct {
	Token        string
	RefreshToken string
	ExpireIn     int64
}

func NewCognito(log *zap.SugaredLogger, sess *session.Session, clientID, userPoolID string) *Cognito {
	identityProvider := cognitoidentityprovider.New(sess)
	return &Cognito{
		log: log, identityProvider: identityProvider,
		clientID:   clientID,
		userPoolID: userPoolID,
	}
}

//CreateUser create a new user inside the cognito pool.
//sub is the unique username of the user, the sub should be unique
//within the pool.
//skipPhoneCheck indicate if we should verify the provided phone number
//by sending sms code or if we should skip the verification and make the
//account active right away
func (c *Cognito) CreateUser(email, phoneNumber, password, aggregator string, skipPhoneCheck bool) (string, error) {
	sub, err := c.GenerateSub(email, phoneNumber, aggregator)

	if err != nil {
		return "", fmt.Errorf("error generating user sub: %v", err)
	}

	inp := cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
			{
				Name:  aws.String("aggregator"),
				Value: aws.String(aggregator),
			},
			{
				Name:  aws.String("isActive"),
				Value: aws.String(strconv.FormatBool(skipPhoneCheck)),
			},
		},
		Username: aws.String(sub),
	}

	if _, err := c.identityProvider.SignUp(&inp); err != nil {
		return "", err
	}

	return sub, nil
}

//ConfirmSignUp validate a newly create account
func (c *Cognito) ConfirmSignUp(code, sub string) error {
	inp := cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		ConfirmationCode: aws.String(code),
		Username:         aws.String(sub),
	}

	if _, err := c.identityProvider.ConfirmSignUp(&inp); err != nil {
		return err
	}

	//Update the user isActive attribute to true
	attr := []*cognitoidentityprovider.AttributeType{
		{
			Name:  aws.String("isActive"),
			Value: aws.String("true"),
		},
	}

	if err := c.updateUserAttribute(sub, attr); err != nil {
		return err
	}

	return nil
}

//AuthenticateUser authenticate a new user, and return
//identification data.
func (c *Cognito) AuthenticateUser(sub, password string) (Session, error) {
	inp := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(sub),
			"SRP_A":    aws.String(password),
		},
		ClientId: aws.String(c.clientID),
	}

	out, err := c.identityProvider.InitiateAuth(&inp)

	if err != nil {
		return Session{}, err
	}

	return Session{
		Token:        *out.AuthenticationResult.AccessToken,
		RefreshToken: *out.AuthenticationResult.RefreshToken,
		ExpireIn:     *out.AuthenticationResult.ExpiresIn,
	}, nil
}

func (c *Cognito) ForgotPassword(sub string) error {
	inp := cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(sub),
	}

	if _, err := c.identityProvider.ForgotPassword(&inp); err != nil {
		return err
	}

	return nil
}

func (c *Cognito) RefreshToken(token string) (Session, error) {
	inp := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeRefreshToken),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(token),
		},
		ClientId: aws.String(c.clientID),
	}

	out, err := c.identityProvider.InitiateAuth(&inp)

	if err != nil {
		return Session{}, err
	}

	return Session{
		Token:        *out.AuthenticationResult.AccessToken,
		RefreshToken: *out.AuthenticationResult.RefreshToken,
		ExpireIn:     *out.AuthenticationResult.ExpiresIn,
	}, nil
}

//DeleteUser will completely delete the user from the pool
//the user will not be able to recover the account.
func (c *Cognito) DeleteUser(accessToken string) error {
	inp := cognitoidentityprovider.DeleteUserInput{
		AccessToken: aws.String(accessToken),
	}

	if _, err := c.identityProvider.DeleteUser(&inp); err != nil {
		return err
	}

	return nil
}

func (c *Cognito) ResendValidateCode(sub string) error {
	inp := cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(sub),
	}

	if _, err := c.identityProvider.ResendConfirmationCode(&inp); err != nil {
		return err
	}

	return nil
}

func (c *Cognito) ConfirmNewPassword(confirmationCode, newPassword, sub string) error {
	inp := cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(c.clientID),
		ConfirmationCode: aws.String(confirmationCode),
		Password:         aws.String(newPassword),
		Username:         aws.String(sub),
	}

	if _, err := c.identityProvider.ConfirmForgotPassword(&inp); err != nil {
		return err
	}

	return nil
}

func (c *Cognito) updateUserAttribute(sub string, attr []*cognitoidentityprovider.AttributeType) error {
	inp := cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: attr,
		UserPoolId:     aws.String(c.userPoolID),
		Username:       aws.String(sub),
	}

	if _, err := c.identityProvider.AdminUpdateUserAttributes(&inp); err != nil {
		return err
	}
	return nil
}

func (c *Cognito) GenerateSub(email, phoneNumber, aggregator string) (string, error) {
	sub := email + phoneNumber + aggregator
	id, err := uuid.FromBytes([]byte(sub))

	if err != nil {
		return "", err
	}

	return id.String(), nil
}
