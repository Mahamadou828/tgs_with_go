package middleware

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {

			v, err := web.GetRequestTrace(ctx)

			log.Infow("request started", "traceId", v.ID, "method", r.Method, "path", r.URL, "remote", r.RemoteAddr, "data", r.Body)

			rqsErr := handler(ctx, w, r)

			if err != nil {
				return rqsErr
			}

			log.Infow("request finished", "traceId", v.ID, "method", r.Method, "path", r.URL, "remote", r.RemoteAddr, "status", v.StatusCode, "since", time.Since(v.Now))

			return nil
		}
		return h
	}

	return m
}
