//Package user provide an api to interact with the business
//logic related to user CRUD and management
package user

import (
	"context"
	"fmt"
	userdto "github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/userroutes/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/user"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log       *zap.SugaredLogger
	db        *sqlx.DB
	aws       *aws.AWS
	userStore user.Store
	aggStore  aggregator.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, aws *aws.AWS) Core {
	return Core{
		log:       log,
		db:        db,
		aws:       aws,
		userStore: user.NewStore(log, db, aws),
		aggStore:  aggregator.NewStore(log, db, aws),
	}
}

type Credentials struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	User         user.User
}

func (c Core) Login(ctx context.Context, aggregator string, payload userdto.Login) (Credentials, error) {
	u, err := c.userStore.QueryByCognitoID(ctx, aggregator, payload.Email, payload.Password)
	if err != nil {
		return Credentials{}, err
	}
	if !u.Active {
		return Credentials{}, fmt.Errorf("user %s is not active", u.Email)
	}

	sess, err := c.aws.Cognito.AuthenticateUser(u.CognitoID, payload.Password)
	if err != nil {
		return Credentials{}, fmt.Errorf("can't authenticate user %s, reason: %v", u.Email, err)
	}

	cred := Credentials{
		Token:        sess.Token,
		RefreshToken: sess.RefreshToken,
		ExpiresIn:    sess.ExpireIn,
		User:         u,
	}

	return cred, nil
}

func (c Core) ConfirmNewPassword(ctx context.Context, payload userdto.ConfirmNewPassword) error {
	u, err := c.userStore.QueryById(ctx, payload.ID)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ConfirmNewPassword(payload.Code, payload.NewPassword, u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) ForgotPassword(ctx context.Context, id string) error {
	u, err := c.userStore.QueryById(ctx, id)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ForgotPassword(u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) VerifyConfirmationCode(ctx context.Context, payload userdto.VerifyConfirmationCode) error {
	u, err := c.userStore.QueryById(ctx, payload.ID)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ConfirmSignUp(payload.ID, u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) ResendConfirmationCode(ctx context.Context, id string) error {
	u, err := c.userStore.QueryById(ctx, id)
	if err != nil {
		return err
	}

	return c.aws.Cognito.ResendValidateCode(u.CognitoID)
}

func (c Core) RefreshToken(ctx context.Context, aggregator string, payload userdto.RefreshToken) (Credentials, error) {
	u, err := c.userStore.QueryById(ctx, payload.ID)
	if err != nil {
		return Credentials{}, err
	}
	if u.AggregatorID != aggregator {
		return Credentials{}, fmt.Errorf("invalid refresh token")
	}

	sess, err := c.aws.Cognito.RefreshToken(payload.Token)
	if err != nil {
		return Credentials{}, fmt.Errorf("session expired")
	}
	cred := Credentials{
		Token:        sess.Token,
		RefreshToken: sess.RefreshToken,
		ExpiresIn:    sess.ExpireIn,
		User:         u,
	}
	return cred, nil
}

func (c Core) Create(ctx context.Context, aggregatorCode string, nu userdto.NewUser, now time.Time) (user.User, error) {
	agg, err := c.aggStore.QueryByID(ctx, aggregatorCode)

	if err != nil {
		return user.User{}, err
	}

	usr, err := c.userStore.Create(ctx, agg.ID, agg.ApiKey, nu, now)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Query(ctx context.Context, pages, rows int) ([]user.User, error) {
	usr, err := c.userStore.Query(ctx, pages, rows)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (user.User, error) {
	usr, err := c.userStore.QueryById(ctx, id)

	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Update(ctx context.Context, id string, ua userdto.UpdateUser, now time.Time) (user.User, error) {
	usr, err := c.userStore.QueryById(ctx, id)

	if err != nil {
		return user.User{}, err
	}

	if ua.Name != nil {
		usr.Name = *ua.Name
	}
	if ua.Email != nil {
		usr.Email = *ua.Email
	}
	if ua.PhoneNumber != nil {
		usr.PhoneNumber = *ua.PhoneNumber
	}
	if ua.Active != nil {
		usr.Active = *ua.Active
	}
	if ua.IsMonthlyActive != nil {
		usr.IsMonthlyActive = *ua.IsMonthlyActive
	}
	if ua.Role != nil {
		usr.Role = *ua.Role
	}
	if ua.IsCGUAccepted != nil {
		usr.IsCGUAccepted = *ua.IsCGUAccepted
	}

	if err := c.userStore.Update(ctx, id, usr, now); err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (c Core) Delete(ctx context.Context, userId string, now time.Time) (user.User, error) {
	u, err := c.userStore.QueryById(ctx, userId)
	if err != nil {
		return user.User{}, err
	}

	if err := c.aws.Cognito.DeleteUser(u.CognitoID); err != nil {
		return user.User{}, err
	}

	usr, err := c.userStore.Delete(ctx, userId, now)
	if err != nil {
		return user.User{}, err
	}

	return usr, nil
}
