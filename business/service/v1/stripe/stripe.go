package stripe

import (
	"errors"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
)

var (
	ErrPaymentAlreadyCancelled = errors.New("payment already canceled")
	ErrPaymentAlreadyRefunded  = errors.New("payment already refunded")
)

type Stripe struct {
}

type PaymentMethod struct {
	Number      string
	CVC         string
	ExpireMonth string
	ExpireYear  string
}

type Payment struct {
	Challenge bool
	Status    stripe.PaymentIntentStatus
	ReturnURL string
}

func New(strID string) Stripe {
	return Stripe{}
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

func CreatePaymentMethod(id string, pm PaymentMethod) (string, error) {
	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			CVC:      stripe.String(pm.CVC),
			ExpMonth: stripe.String(pm.ExpireMonth),
			ExpYear:  stripe.String(pm.ExpireYear),
			Number:   stripe.String(pm.Number),
		},
		Customer: stripe.String(id),
	}

	res, err := paymentmethod.New(params)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func CreateCharge(strKey string, cusID, pmID string, amount int64, returnURL string) (Payment, error) {
	stripe.Key = strKey

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amount),
		Confirm:       stripe.Bool(true),
		Currency:      stripe.String("eur"),
		Customer:      stripe.String(cusID),
		PaymentMethod: stripe.String(pmID),
		ReturnURL:     stripe.String(returnURL),
	}

	res, err := paymentintent.New(params)
	if err != nil {
		return Payment{}, err
	}

	if res.NextAction != nil {
		return Payment{true, res.Status, res.NextAction.RedirectToURL.URL}, nil
	}

	return Payment{false, res.Status, ""}, nil
}

func isPaymentValidate(strKey string, payID string) (bool, error) {
	stripe.Key = strKey

	pi, err := paymentintent.Get(payID, nil)
	if err != nil {
		return false, err
	}

	return pi.Status == stripe.PaymentIntentStatusSucceeded, nil
}

func CancelPayment(strKey string, id string) error {
	stripe.Key = strKey

	if _, err := paymentintent.Cancel(id, nil); err != nil {
		return err
	}
	return nil
}

func RefundUser(strKey string, piID string, amount int64) error {
	stripe.Key = strKey

	pi, err := paymentintent.Get(piID, nil)
	if err != nil {
		return err
	}

	if pi.Status == stripe.PaymentIntentStatusCanceled {
		return ErrPaymentAlreadyCancelled
	}

	params := &stripe.RefundParams{
		Amount:        stripe.Int64(amount),
		PaymentIntent: stripe.String(piID),
	}

	if _, err := refund.New(params); err != nil {
		return err
	}
	return nil
}
