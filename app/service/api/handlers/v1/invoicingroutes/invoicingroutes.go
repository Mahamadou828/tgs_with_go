package invoicingroutes

import (
	"context"
	"fmt"
	invoicingentity2 "github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/invoicingentity"
	"net/http"
	"strconv"

	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/invoicingentity"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
)

type Handler struct {
	InCore invoicingentity.Core
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

	es, err := h.InCore.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(
		ctx,
		w,
		http.StatusOK,
		es,
	)
}

func (h Handler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	e, err := h.InCore.QueryByID(ctx, id)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, e)
}

func (h Handler) QueryByEnterprise(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	c := web.Param(r, "id")

	e, err := h.InCore.QueryByEnterpriseID(ctx, c)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, e)
}

func (h Handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)
	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	var eu invoicingentity2.NewInvoicingEntityDTO

	if err := web.Decode(r, &eu); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(eu); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	e, err := h.InCore.Create(ctx, eu, v.Now)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, e)
}

func (h Handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	v, err := web.GetRequestTrace(ctx)
	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	var eu invoicingentity2.UpdateInvoicingEntityDTO

	if err := web.Decode(r, &eu); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(eu); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	e, err := h.InCore.Update(ctx, id, eu, v.Now)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, e)
}

func (h Handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	v, err := web.GetRequestTrace(ctx)
	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	e, err := h.InCore.Delete(ctx, id, v.Now)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, e)
}
