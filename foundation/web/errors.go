package web

import "fmt"

type RequestError struct {
	Message error
	Details []string
	Status  int
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int, details ...string) *RequestError {
	return &RequestError{
		Message: err,
		Status:  status,
		Details: details,
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (e *RequestError) Error() string {
	err := fmt.Errorf("error: %w, details: %v", e.Message, e.Details)
	return err.Error()
}
