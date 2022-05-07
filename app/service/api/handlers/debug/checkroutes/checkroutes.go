//Package checkroutes contain all route for the debug server
package checkroutes

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

/**
@todo make readiness check id the database support is up
*/

import (
	"encoding/json"
)

type Handler struct {
	Build  string
	Logger *zap.SugaredLogger
	DB     *sqlx.DB
}

// Readiness checks if the service is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	if err := database.StatusCheck(ctx, h.DB); err != nil {
		status = "not ok"
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)

	if err != nil {
		h.Logger.Errorw("readiness", "Error", err)
	}

	if _, err := w.Write(jsonData); err != nil {
		h.Logger.Errorw("readiness", "Error", err)
	}
}

// Liveliness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h Handler) Liveliness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	statusCode := http.StatusOK

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)

	if err != nil {
		h.Logger.Errorw("readiness", "Error", err)
	}

	if _, err := w.Write(jsonData); err != nil {
		h.Logger.Errorw("readiness", "Error", err)
	}
}
