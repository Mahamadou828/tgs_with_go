package middleware

import (
	"context"
	"fmt"
	v1 "github.com/Mahamadou828/tgs_with_golang/business/web/v1"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"go.uber.org/zap"
	"net/http"
)

//Errors send error handler error to the client after formatting them into the ErrorResponse.
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
			//Get the request trace from the context
			v, err := web.GetRequestTrace(ctx)

			//If there is no trace, we should shutdown the
			//app because we have an integrity issue
			if err != nil {
				panic(err)
			}

			rqsErr := handler(ctx, w, r)

			if rqsErr != nil {
				//Log the error
				log.Errorw("request error", "traceId", v.ID, "method", r.Method, "path", r.URL, "error", rqsErr.Message, "details", rqsErr.Details)
				rsp := v1.ErrorResponse{
					Message: rqsErr.Message.Error(),
					Details: rqsErr.Details,
					TraceID: v.ID,
				}

				if err := web.Response(ctx, w, rqsErr.Status, rsp); err != nil {
					return &web.RequestError{
						Message: fmt.Errorf("can't send http response: %v", err),
						Details: nil,
						Status:  http.StatusInternalServerError,
					}
				}
			}

			return nil
		}

		return h
	}

	return m
}
