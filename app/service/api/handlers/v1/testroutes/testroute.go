package testroutes

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Logger *zap.SugaredLogger
	Build  string
}

func (h Handler) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "Ok",
	}

	if err := web.Response(ctx, w, http.StatusOK, resp); err != nil {
		return web.NewRequestError(fmt.Errorf("can't send response: %v", err), http.StatusInternalServerError)
	}

	return nil
}

func (h Handler) TestFail(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	return &web.RequestError{
		Message: fmt.Errorf("invalid formular"),
		Details: []string{"Field ID must not be empty", "Email is not an email", "phone number is not valid"},
		Status:  http.StatusBadRequest,
	}
}
