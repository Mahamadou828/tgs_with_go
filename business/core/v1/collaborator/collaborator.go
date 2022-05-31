package collaborator

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/collaborator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterprisepolicy"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/enterpriseteam"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

//All collaborator are associated with the same aggregator
const aggregatorCode = "tgs-corporate"

type Core struct {
	collaboratorStore collaborator.Store
	aws               *aws.AWS
	db                *sqlx.DB
	log               *zap.SugaredLogger
	aggregatorStore   aggregator.Store
	policyStore       enterprisepolicy.Store
	teamStore         enterpriseteam.Store
}

type Credentials struct {
	Token        string                    `json:"token"`
	RefreshToken string                    `json:"refreshToken"`
	ExpiresIn    int64                     `json:"expiresIn"`
	Collaborator collaborator.Collaborator `json:"collaborator"`
}

func NewCore(aws *aws.AWS, db *sqlx.DB, log *zap.SugaredLogger) Core {
	return Core{
		aws:               aws,
		db:                db,
		log:               log,
		collaboratorStore: collaborator.NewStore(log, db, aws),
		aggregatorStore:   aggregator.NewStore(log, db, aws),
		policyStore:       enterprisepolicy.NewStore(db, log),
		teamStore:         enterpriseteam.NewStore(db, log),
	}
}

func (c Core) Login(ctx context.Context, payload dto.Login) (Credentials, error) {
	agg, err := c.aggregatorStore.QueryByCode(ctx, aggregatorCode)
	//if the tgs-corporate aggregator does not exist we should panic
	//because we have an integrity issue
	if err != nil {
		panic(err)
	}
	co, err := c.collaboratorStore.QueryByEmail(ctx, agg.ID, payload.Email)
	if err != nil {
		return Credentials{}, err
	}
	if !co.Active {
		return Credentials{}, fmt.Errorf("user: %s is not active", payload.Email)
	}

	sess, err := c.aws.Cognito.AuthenticateUser(co.CognitoID, payload.Password)
	if err != nil {
		return Credentials{}, err
	}
	cred := Credentials{
		Token:        sess.Token,
		RefreshToken: sess.RefreshToken,
		ExpiresIn:    sess.ExpireIn,
		Collaborator: co,
	}

	return cred, nil
}

func (c Core) RefreshToken(ctx context.Context, payload dto.RefreshToken) (Credentials, error) {
	co, err := c.collaboratorStore.QueryByID(ctx, payload.ID)
	if err != nil {
		return Credentials{}, err
	}

	sess, err := c.aws.Cognito.RefreshToken(payload.RefreshToken)
	if err != nil {
		return Credentials{}, err
	}
	cred := Credentials{
		Token:        sess.Token,
		RefreshToken: sess.RefreshToken,
		ExpiresIn:    sess.ExpireIn,
		Collaborator: co,
	}
	return cred, nil
}

func (c Core) Create(ctx context.Context, nco dto.NewCollaborator, now time.Time) (collaborator.Collaborator, error) {
	agg, err := c.aggregatorStore.QueryByCode(ctx, aggregatorCode)
	//if the tgs-corporate aggregator does not exist we should panic
	//because we have an integrity issue
	if err != nil {
		panic(err)
	}

	t, err := c.teamStore.QueryByID(ctx, nco.TeamID)
	if err != nil {
		return collaborator.Collaborator{}, err
	}

	p, err := c.policyStore.QueryByID(ctx, t.PolicyID)
	if err != nil {
		return collaborator.Collaborator{}, err
	}

	co, err := c.collaboratorStore.Create(
		ctx,
		agg,
		nco,
		p.CollaboratorBudget,
		now,
	)

	if err != nil {
		return collaborator.Collaborator{}, err
	}

	return co, nil
}

func (c Core) QueryByID(ctx context.Context, id string) (collaborator.Collaborator, error) {
	co, err := c.collaboratorStore.QueryByID(ctx, id)
	if err != nil {
		return collaborator.Collaborator{}, err
	}
	return co, nil
}

func (c Core) Query(ctx context.Context, pageNumber, rowsPerPage int) ([]collaborator.Collaborator, error) {
	cos, err := c.collaboratorStore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return []collaborator.Collaborator{}, err
	}
	return cos, nil
}

func (c Core) QueryByEnterprise(ctx context.Context, id string, pageNumber, rowsPerPage int) ([]collaborator.Collaborator, error) {
	cos, err := c.collaboratorStore.QueryByEnterprise(ctx, id, pageNumber, rowsPerPage)
	if err != nil {
		return []collaborator.Collaborator{}, err
	}
	return cos, nil
}

func (c Core) Update(ctx context.Context, id string, uc dto.UpdateCollaborator, now time.Time) (collaborator.Collaborator, error) {
	co, err := c.collaboratorStore.QueryByID(ctx, id)

	if err != nil {
		return collaborator.Collaborator{}, err
	}

	if uc.Name != nil {
		co.Name = *uc.Name
	}
	if uc.Email != nil {
		co.Email = *uc.Email
	}
	if uc.PhoneNumber != nil {
		co.PhoneNumber = *uc.PhoneNumber
	}
	if uc.Active != nil {
		co.Active = *uc.Active
	}
	if uc.IsMonthlyActive != nil {
		co.IsMonthlyActive = *uc.IsMonthlyActive
	}
	if uc.Role != nil {
		co.Role = *uc.Role
	}
	if uc.IsCGUAccepted != nil {
		co.IsCGUAccepted = *uc.IsCGUAccepted
	}

	if err := c.collaboratorStore.Update(ctx, id, co, now); err != nil {
		return collaborator.Collaborator{}, err
	}

	return co, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (collaborator.Collaborator, error) {
	co, err := c.collaboratorStore.QueryByID(ctx, id)
	if err != nil {
		return collaborator.Collaborator{}, err
	}
	if err := c.collaboratorStore.Delete(ctx, id, now); err != nil {
		return collaborator.Collaborator{}, err
	}
	return co, nil
}

func (c Core) ConfirmNewPassword(ctx context.Context, payload dto.ConfirmNewPassword) error {
	u, err := c.collaboratorStore.QueryByID(ctx, payload.ID)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ConfirmNewPassword(payload.Code, payload.NewPassword, u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) ForgotPassword(ctx context.Context, id string) error {
	u, err := c.collaboratorStore.QueryByID(ctx, id)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ForgotPassword(u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) VerifyConfirmationCode(ctx context.Context, payload dto.VerifyConfirmationCode) error {
	u, err := c.collaboratorStore.QueryByID(ctx, payload.ID)
	if err != nil {
		return err
	}
	if err := c.aws.Cognito.ConfirmSignUp(payload.ID, u.CognitoID); err != nil {
		return err
	}
	return nil
}

func (c Core) ResendConfirmationCode(ctx context.Context, id string) error {
	u, err := c.collaboratorStore.QueryByID(ctx, id)
	if err != nil {
		return err
	}

	return c.aws.Cognito.ResendValidateCode(u.CognitoID)
}
