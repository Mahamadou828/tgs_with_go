package aggregatorroutes

import (
	"context"
	"fmt"
	"net/http"

	aggCore "github.com/Mahamadou828/tgs_with_golang/business/core/v1/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
)

type Handler struct {
	Agg aggCore.Core
}

func (h Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)

	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	var na aggregator.NewAggregator

	if err := web.Decode(r, &na); err != nil {
		return web.NewRequestError(
			fmt.Errorf("unable to decode payload: %v", err),
			http.StatusInternalServerError,
		)
	}

	agr, err := h.Agg.Create(ctx, na, v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusCreated, agr)
}
