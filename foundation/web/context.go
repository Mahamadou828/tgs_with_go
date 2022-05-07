package web

import (
	"context"
	"errors"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is how request values are stored/retrieved.
const key ctxKey = 1

//RequestTrace is a unique value attach to each request
//It's contains information about the request context
//and it's use in logging tracing and sentry output
type RequestTrace struct {
	ID         string
	Now        time.Time
	StatusCode int
}

// GetRequestTrace returns the values from the context.
func GetRequestTrace(ctx context.Context) (*RequestTrace, error) {
	v, ok := ctx.Value(key).(*RequestTrace)

	if !ok {
		return nil, errors.New("web request trace not found")
	}

	return v, nil
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*RequestTrace)

	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.ID
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*RequestTrace)
	if !ok {
		return errors.New("web value missing from context")
	}
	v.StatusCode = statusCode
	return nil
}
