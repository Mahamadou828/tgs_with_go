package userroutes

import (
	"context"
	"fmt"
	aggCore "github.com/Mahamadou828/tgs_with_golang/business/core/v1/aggregator"
	userCore "github.com/Mahamadou828/tgs_with_golang/business/core/v1/user"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/store/user"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
	"strconv"
)

type Handler struct {
	User userCore.Core
	Agg  aggCore.Core
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
		return web.NewRequestError(fmt.Errorf("invalid rows format [%s]", rows), http.StatusBadRequest)
	}

	u, err := h.User.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, u)
}

func (h Handler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	u, err := h.User.QueryByID(ctx, id)

	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, u)
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

	if err := web.Decode(r, &nu); err != nil {
		return web.NewRequestError(
			fmt.Errorf("unable to decode payload: %v", err),
			http.StatusInternalServerError,
		)
	}

	usr, err := h.User.Create(ctx, r.Header.Get("aggregator"), nu, v.Now)

	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("can't create user: %q, reason: %v", nu.Email, err),
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusCreated, usr)
}

func (h Handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	v, err := web.GetRequestTrace(ctx)

	if err != nil {
		return web.NewRequestError(
			web.NewShutdownError("web value missing from context"),
			http.StatusInternalServerError,
		)
	}

	u, err := h.User.Delete(ctx, web.Param(r, "id"), v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, u)
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

	var ua user.UpdateUser
	if err := web.Decode(r, &ua); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(ua); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	u, err := h.User.Update(ctx, id, ua, v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, u)
}
