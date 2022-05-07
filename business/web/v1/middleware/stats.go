package middleware

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
)

func Stats() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
			//@todo have access to the database and save the current request response.
			//Maybe even put that logic inside the error middleware
			return nil
		}
		return h
	}
	return m
}