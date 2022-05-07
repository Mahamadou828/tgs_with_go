package middleware

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/metrics"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
	"runtime/debug"
)

func Panic() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err *web.RequestError) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					// Stack trace will be provided.
					err = web.NewRequestError(fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace)), http.StatusInternalServerError)

					metrics.AddPanics(ctx)
				}
			}()
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
