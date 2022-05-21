//Package handlers manage the different version of the api
package handlers

import (
	"expvar"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/aggregator"
	"net/http"
	"net/http/pprof"

	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/debug/checkroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/aggregatorroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/testroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/userroutes"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/user"
	"github.com/Mahamadou828/tgs_with_golang/business/web/v1/middleware"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func APIMux(cfg web.AppConfig) *web.App {
	const version = "v1"
	//Create a new app instance
	app := web.NewApp(
		cfg,
		version,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Cors(cfg.CorsOrigin),
		middleware.Panic(),
	)

	//Load the v1 route
	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg web.AppConfig) {
	trt := testroutes.Handler{
		Logger: cfg.Log,
		Build:  cfg.Build,
		Env:    cfg.Env,
	}

	urt := userroutes.Handler{User: user.NewCore(cfg.Log, cfg.DB, cfg.AWS)}
	agt := aggregatorroutes.Handler{Agg: aggregator.NewCore(cfg.Log, cfg.DB, cfg.AWS)}

	//=========================== Test Route
	app.Handle(http.MethodGet, "/test", trt.Test)
	app.Handle(http.MethodGet, "/test/fail", trt.TestFail)
	app.Handle(http.MethodGet, "/test/panic", trt.TestPanic)

	//=========================== User Route
	app.Handle(http.MethodPost, "/user", urt.Create)
	app.Handle(http.MethodGet, "/user/:id", urt.QueryByID)
	app.Handle(http.MethodGet, "/user", urt.Query)
	app.Handle(http.MethodPut, "/user/:id", urt.Update)
	app.Handle(http.MethodDelete, "/user/:id", urt.Delete)
	//=========================== User Auth Route
	app.Handle(http.MethodPost, "/user/login", urt.Login)
	app.Handle(http.MethodPost, "/user/token/refresh", urt.RefreshToken)
	app.Handle(http.MethodPost, "/user/code/verify", urt.VerifyConfirmationCode)
	app.Handle(http.MethodGet, "/user/code/resend/:id", urt.ResendConfirmationCode)
	app.Handle(http.MethodPut, "/user/password/forgot/:id", urt.ForgotPassword)
	app.Handle(http.MethodPut, "/user/password/reset", urt.ConfirmNewPassword)

	//=========================== Aggregator Route
	app.Handle(http.MethodPost, "/aggregator", agt.Create)
	app.Handle(http.MethodPut, "/aggregator/:id", agt.Update)
	app.Handle(http.MethodDelete, "/aggregator/:id", agt.Delete)
	app.Handle(http.MethodGet, "/aggregator", agt.Query)
	app.Handle(http.MethodGet, "/aggregator/:id", agt.QueryByID)
}

func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) *http.ServeMux {
	mux := DebugStandardLibraryMux()

	handlers := checkroutes.Handler{
		Build:  build,
		Logger: log,
		DB:     db,
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
