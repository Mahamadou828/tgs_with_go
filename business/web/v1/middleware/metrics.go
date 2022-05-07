package middleware

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/metrics"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
)

//Metrics updates program counters
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
			//Add the metrics into the context for metrics gathering
			metrics.Set(ctx)
			//Call the next handler
			err := handler(ctx, w, r)
			//Handle updating metrics that can be handled here
			//Increment the number of requests and goroutines
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)
			//Increment if there is an error flowing through the request
			if err != nil {
				metrics.AddErrors(ctx)
			}
			//Return the error so it can be handled properly
			return err
		}

		return h
	}

	return m
}
