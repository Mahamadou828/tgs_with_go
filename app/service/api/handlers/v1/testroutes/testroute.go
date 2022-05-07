package testroutes

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Logger *zap.SugaredLogger
	Build  string
}

func (h Handler) Test(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "Ok",
	}

	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(resp)

	if err != nil {
		h.Logger.Errorw("/api/test", "status", "test route error", "error", err)
	}

	if _, err := w.Write(jsonData); err != nil {
		h.Logger.Errorw("/api/test", "status", "test route error", "error", err)
	}
}
