package middleware

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
)

func Cors(origin string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
			// Set the CORS headers to the response.
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, aggregatorCode, x-api-token, enterpriseCode")

			return handler(ctx, w, r)
		}
		return h
	}

	return m
}