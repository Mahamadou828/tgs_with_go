package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/setupintent"
)

type PaymentMethod struct {
	ID                   string
	IsThreeDSecureNeeded bool
	ThreeDSecureURL      string
}

type PaymentMethodParams struct {
	Number      string
	CVC         string
	ExpireMonth string
	ExpireYear  string
	ReturnURL   string
}

type CreateChargeParams struct {
	CusID     string
	PmID      string
	Amount    int64
	ReturnURL string
	Currency  string
}

type Charge struct {
	Challenge bool
	Status    stripe.PaymentIntentStatus
	ReturnURL string
}

//CreateUser creates a new stripe customer and return his id
func CreateUser(strKey, email, phoneNumber, name string) (string, error) {
	stripe.Key = strKey

	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Phone: stripe.String(phoneNumber),
	}

	c, err := customer.New(params)
	if err != nil {
		return "", err
	}

	return c.ID, nil
}

func DeleteUser(strKey string, id string) error {
	stripe.Key = strKey

	if _, err := customer.Del(id, nil); err != nil {
		return err
	}
	return nil
}

func CreatePaymentMethod(strKey string, cusID string, pm PaymentMethodParams) (PaymentMethod, error) {
	stripe.Key = strKey

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			CVC:      stripe.String(pm.CVC),
			ExpMonth: stripe.String(pm.ExpireMonth),
			ExpYear:  stripe.String(pm.ExpireYear),
			Number:   stripe.String(pm.Number),
		},
		Type: stripe.String("card"),
	}

	p, err := paymentmethod.New(params)
	if err != nil {
		return PaymentMethod{}, err
	}

	i, err := setupintent.New(&stripe.SetupIntentParams{
		Customer:           stripe.String(cusID),
		PaymentMethodTypes: []*string{stripe.String("card")},
		Confirm:            stripe.Bool(true),
		PaymentMethod:      stripe.String(p.ID),
	})
	if err != nil {
		return PaymentMethod{}, err
	}

	if i.NextAction == nil {
		return PaymentMethod{p.ID, false, ""}, nil
	}

	return PaymentMethod{p.ID, true, i.NextAction.RedirectToURL.URL}, nil
}

func CreateCharge(strKey string, p CreateChargeParams) (Charge, error) {
	stripe.Key = strKey

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(p.Amount),
		Confirm:       stripe.Bool(true),
		Currency:      stripe.String(p.Currency),
		Customer:      stripe.String(p.CusID),
		PaymentMethod: stripe.String(p.PmID),
		ReturnURL:     stripe.String(p.ReturnURL),
		OffSession:    stripe.Bool(true),
		CaptureMethod: stripe.String("manual"),
	}

	res, err := paymentintent.New(params)
	if err != nil {
		return Charge{}, err
	}

	if res.NextAction != nil {
		return Charge{true, res.Status, res.NextAction.RedirectToURL.URL}, nil
	}

	return Charge{false, res.Status, ""}, nil
}

func CancelCharge(strKey string, id string) error {
	stripe.Key = strKey

	if _, err := paymentintent.Cancel(id, nil); err != nil {
		return err
	}
	return nil
}

func CaptureCharge(strKey string, ptID string, amount int64) error {
	stripe.Key = strKey

	_, err := paymentintent.Capture(ptID, &stripe.PaymentIntentCaptureParams{
		AmountToCapture: stripe.Int64(amount),
	})
	if err != nil {
		return err
	}
	return nil
}
