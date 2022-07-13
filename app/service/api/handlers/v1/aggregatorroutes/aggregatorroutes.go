package aggregatorroutes

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/aggregator"
	"net/http"
	"strconv"

	aggCore "github.com/Mahamadou828/tgs_with_golang/business/core/v1/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
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

	var na aggregator.NewAggregatorDTO

	if err := web.Decode(r, &na); err != nil {
		return web.NewRequestError(
			fmt.Errorf("unable to decode payload: %v", err),
			http.StatusInternalServerError,
		)
	}

	if err := validate.Check(na); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
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

func (h Handler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")

	agg, err := h.Agg.QueryByID(ctx, id)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, agg)
}

func (h Handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)
	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	aggr, err := h.Agg.Delete(ctx, web.Param(r, "id"), v.Now)

	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, aggr)
}

func (h Handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)
	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	var agg aggregator.UpdateAggregatorDTO
	if err := web.Decode(r, &agg); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(agg); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	aggr, err := h.Agg.Update(ctx, id, agg, v.Now)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, aggr)
}

func (h Handler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	page := web.QueryParam(r, "page")
	if len(page) != 1 {
		return web.NewRequestError(fmt.Errorf("missing query parameter page"), http.StatusBadRequest)
	}
	pageNumber, err := strconv.Atoi(page[0])
	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid page format [%s]", page),
			http.StatusBadRequest,
		)
	}

	rows := web.QueryParam(r, "rows")
	if len(page) != 1 {
		return web.NewRequestError(fmt.Errorf("missing query parameter rows"), http.StatusBadRequest)
	}
	rowsPerPage, err := strconv.Atoi(rows[0])
	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid rows format [%s]", rows),
			http.StatusBadRequest,
		)
	}

	aggrs, err := h.Agg.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, aggrs)
}
