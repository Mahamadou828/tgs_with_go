package testroutes

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
)

type Handler struct {
	Logger *zap.SugaredLogger
	Build  string
	Env    string
}

func (h Handler) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	resp := struct {
		Status string `json:"status"`
		Build  string `json:"build"`
		Env    string `json:"env"`
	}{
		Status: "Ok",
		Build:  h.Build,
		Env:    h.Env,
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

func (h Handler) TestPanic(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	err := fmt.Errorf("panic")
	panic(err)
	return nil
}
