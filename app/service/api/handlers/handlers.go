//Package handlers manage the different version of the api
package handlers

import (
	"expvar"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/debug/checkroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/testroutes"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
)

func ApiMux(build string, log *zap.SugaredLogger) *http.ServeMux {
	mux := http.NewServeMux()

	handlers := testroutes.Handler{
		Logger: log,
		Build:  build,
	}

	mux.HandleFunc("/api/test", handlers.Test)

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) *http.ServeMux {
	mux := DebugStandardLibraryMux()

	handlers := checkroutes.Handler{
		Build:  build,
		Logger: log,
	}

	mux.HandleFunc("/debug/readiness", handlers.Readiness)
	mux.HandleFunc("/debug/liveliness", handlers.Liveliness)

	return mux
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}