package aggregatorroutes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

	var agg aggregator.UpdateAggregator
	if err := web.Decode(r, &agg); err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	aggr, err := h.Agg.Update(ctx, id, agg, v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, aggr)
}

func (h Handler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	page := web.QueryParam(r, "page")[0]
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("invalid page format [%s]", page),
			http.StatusBadRequest,
		)
	}

	rows := web.QueryParam(r, "rows")[0]
	rowsPerPage, err := strconv.Atoi(rows)
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
