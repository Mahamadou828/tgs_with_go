package paymentmethod

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/paymentmethod"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/user"
	"github.com/Mahamadou828/tgs_with_golang/business/service/v1/stripe"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	log       *zap.SugaredLogger
	db        *sqlx.DB
	pmStore   paymentmethod.Store
	userStore user.Store
	stripeKey string
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB, stripeKey string) Core {
	return Core{
		log:       log,
		db:        db,
		pmStore:   paymentmethod.NewStore(log, db),
		userStore: user.NewStore(log, db),
		stripeKey: stripeKey,
	}
}

func (c Core) Query(ctx context.Context, id string, pageNumber, rowsPerPage int) ([]paymentmethod.PaymentMethod, error) {
	pms, err := c.pmStore.Query(ctx, id, pageNumber, rowsPerPage)
	if err != nil {
		return []paymentmethod.PaymentMethod{}, err
	}

	return pms, nil
}

func (c Core) Create(ctx context.Context, npm paymentmethod.NewPaymentMethodDTO, now time.Time) (
	struct {
		IsThreeDSecureNeeded bool                        `json:"isThreeDSecureNeeded"`
		ReturnUrl            string                      `json:"returnUrl"`
		Card                 paymentmethod.PaymentMethod `json:"card"`
	},
	error,
) {
	var res struct {
		IsThreeDSecureNeeded bool                        `json:"isThreeDSecureNeeded"`
		ReturnUrl            string                      `json:"returnUrl"`
		Card                 paymentmethod.PaymentMethod `json:"card"`
	}
	u, err := c.userStore.QueryByID(ctx, npm.UserID)

	if err != nil {
		return res, err
	}

	spm, err := stripe.CreatePaymentMethod(
		c.stripeKey,
		u.StripeID,
		stripe.PaymentMethodParams{
			Number:      npm.Number,
			CVC:         npm.CVC,
			ExpireMonth: npm.ExpireMonth,
			ExpireYear:  npm.ExpireYear,
			ReturnURL:   npm.ReturnURL,
		},
	)
	if err != nil {
		return res, err
	}

	if spm.IsThreeDSecureNeeded {
		res.IsThreeDSecureNeeded, res.ReturnUrl = spm.IsThreeDSecureNeeded, spm.ThreeDSecureURL
	}

	p, err := c.pmStore.Create(ctx, spm.ID, npm, now)
	if err != nil {
		return res, err
	}

	res.Card = p

	return res, nil
}

func (c Core) Update(ctx context.Context, id string, upm paymentmethod.UpdatePaymentMethodDTO, now time.Time) (paymentmethod.PaymentMethod, error) {
	pm, err := c.pmStore.QueryByID(ctx, id)
	if err != nil {
		return paymentmethod.PaymentMethod{}, err
	}

	if upm.IsFavorite != nil {
		pm.IsFavorite = *upm.IsFavorite
	}
	if upm.Name != nil {
		pm.Name = *upm.Name
	}

	if err := c.pmStore.Update(ctx, pm, now); err != nil {
		return paymentmethod.PaymentMethod{}, err
	}

	return pm, nil
}

func (c Core) Delete(ctx context.Context, id string, now time.Time) (paymentmethod.PaymentMethod, error) {
	p, err := c.pmStore.QueryByID(ctx, id)
	if err != nil {
		return paymentmethod.PaymentMethod{}, err
	}

	if err := c.pmStore.Delete(ctx, id, now); err != nil {
		return paymentmethod.PaymentMethod{}, err
	}

	return p, err
}
