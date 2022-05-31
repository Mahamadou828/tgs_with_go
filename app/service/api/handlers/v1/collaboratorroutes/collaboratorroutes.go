package collaboratorroutes

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/collaborator"
	"github.com/Mahamadou828/tgs_with_golang/business/data/v1/dto"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/validate"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
	"strconv"
)

type Handler struct {
	Co collaborator.Core
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

	u, err := h.Co.Query(ctx, pageNumber, rowsPerPage)
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

	u, err := h.Co.QueryByID(ctx, id)

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

	var nc dto.NewCollaborator

	if err := web.Decode(r, &nc); err != nil {
		return web.NewRequestError(
			fmt.Errorf("unable to decode payload: %v", err),
			http.StatusInternalServerError,
		)
	}

	usr, err := h.Co.Create(ctx, nc, v.Now)

	if err != nil {
		return web.NewRequestError(
			fmt.Errorf("can't create user: %q, reason: %v", nc.Email, err),
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

	u, err := h.Co.Delete(ctx, web.Param(r, "id"), v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}
	return web.Response(ctx, w, http.StatusOK, u)
}

func (h Handler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	var login dto.Login

	if err := web.Decode(r, &login); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(login); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	cred, err := h.Co.Login(ctx, login)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, cred)
}

func (h Handler) ResendConfirmationCode(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if err := h.Co.ResendConfirmationCode(ctx, id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusNoContent, nil)
}

func (h Handler) VerifyConfirmationCode(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	var payload dto.VerifyConfirmationCode
	if err := web.Decode(r, &payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := h.Co.VerifyConfirmationCode(ctx, payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	return web.Response(ctx, w, http.StatusNoContent, nil)
}

func (h Handler) ForgotPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	id := web.Param(r, "id")
	if err := validate.CheckID(id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := h.Co.ForgotPassword(ctx, id); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	return web.Response(ctx, w, http.StatusNoContent, nil)
}

func (h Handler) ConfirmNewPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	var payload dto.ConfirmNewPassword
	if err := web.Decode(r, &payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := h.Co.ConfirmNewPassword(ctx, payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	return web.Response(ctx, w, http.StatusNoContent, nil)
}

func (h Handler) RefreshToken(ctx context.Context, w http.ResponseWriter, r *http.Request) *web.RequestError {
	var payload dto.RefreshToken

	if err := web.Decode(r, &payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	cred, err := h.Co.RefreshToken(ctx, payload)
	if err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Response(ctx, w, http.StatusOK, cred)
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

	var ua dto.UpdateCollaborator
	if err := web.Decode(r, &ua); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}
	if err := validate.Check(ua); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	u, err := h.Co.Update(ctx, id, ua, v.Now)

	if err != nil {
		return web.NewRequestError(
			err,
			http.StatusBadRequest,
		)
	}

	return web.Response(ctx, w, http.StatusOK, u)
}
