package middleware

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/metrics"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"github.com/getsentry/sentry-go"
	"net/http"
)

func Panic() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err *web.RequestError) {
			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			//@todo review the panic handling
			defer func() {
				if rec := recover(); rec != nil {
					// Stack trace will be provided.
					v, err := web.GetRequestTrace(ctx)

					//If there's no request trace associated with the current request
					//we should shut down the system because we have an integrity issue
					if err != nil {
						panic(err)
					}

					err = web.NewRequestError(
						fmt.Errorf("panic: %v", sentry.NewRequest(r)),
						http.StatusInternalServerError,
					)
					eventID := v.Hub.RecoverWithContext(context.WithValue(ctx, sentry.RequestContextKey, r), err)

					if eventID != nil {
						v.Hub.Flush(5)
					}

					metrics.AddPanics(ctx)

					if err := web.Response(ctx, w, http.StatusInternalServerError, struct {
						Message string
					}{
						Message: "Internal Server Error",
					}); err != nil {
						panic(err)
					}
				}
			}()
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
