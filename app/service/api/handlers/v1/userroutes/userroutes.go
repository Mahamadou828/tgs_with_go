package userroutes

import (
	"context"
	"fmt"
	userCore "github.com/Mahamadou828/tgs_with_golang/business/core/v1/user"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/user"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
)

type Handler struct {
	User userCore.Core
}

func (h Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)

	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	var nu user.NewUser

	if err := web.Decode(r, nu); err != nil {
		return web.NewRequestError(
			fmt.Errorf("unable to decode payload: %v", err),
			http.StatusInternalServerError,
		)
	}

	aggregatorId, apiKey := r.Header.Get("aggregatorId"), r.Header.Get("apiKey")

	//@todo check if the aggregator exist inside the db
	if aggregatorId == "" || apiKey == "" {
		return web.NewRequestError(
			fmt.Errorf("aggregatorId and apiKey cannot be empty: %q, %q", aggregatorId, apiKey),
			http.StatusBadRequest,
		)
	}

	usr, err := h.User.Create(ctx, aggregatorId, apiKey, nu, v.Now)

	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("can't create user: %q, reason: %v", nu.Email, err),
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusCreated, usr)
}
