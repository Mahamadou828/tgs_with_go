//Package handlers manage the different version of the api
package handlers

import (
	"expvar"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/enterprisepolicyroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/paymentmethodroutes"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/enterprisepolicy"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/enterpriseteam"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/invoicingentity"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/paymentmethod"
	"net/http"
	"net/http/pprof"

	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/debug/checkroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/aggregatorroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/collaboratorroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/enterprisepackroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/enterpriseroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/enterpriseteamroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/invoicingroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/testroutes"
	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers/v1/userroutes"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/aggregator"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/collaborator"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/enterprise"
	"github.com/Mahamadou828/tgs_with_golang/business/core/v1/enterprisepack"
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

	ugt := userroutes.Handler{User: user.NewCore(cfg.Log, cfg.DB, cfg.AWS, cfg.StripeKey)}
	agt := aggregatorroutes.Handler{Agg: aggregator.NewCore(cfg.Log, cfg.DB, cfg.AWS)}
	egt := enterpriseroutes.Handler{En: enterprise.NewCore(cfg.Log, cfg.DB)}
	epgt := enterprisepackroutes.Handler{Pac: enterprisepack.NewCore(cfg.Log, cfg.DB)}
	cgt := collaboratorroutes.Handler{Co: collaborator.NewCore(cfg.AWS, cfg.DB, cfg.Log, cfg.StripeKey)}
	tgt := enterpriseteamroutes.Handler{TeCore: enterpriseteam.NewCore(cfg.DB, cfg.Log)}
	tpt := enterprisepolicyroutes.Handler{PoCore: enterprisepolicy.NewCore(cfg.DB, cfg.Log)}
	igt := invoicingroutes.Handler{InCore: invoicingentity.NewCore(cfg.Log, cfg.DB)}
	pmgt := paymentmethodroutes.Handler{PmCore: paymentmethod.NewCore(cfg.Log, cfg.DB, cfg.StripeKey)}

	//=========================== Test Route
	app.Handle(http.MethodGet, "/test", trt.Test)
	app.Handle(http.MethodGet, "/test/fail", trt.TestFail)
	app.Handle(http.MethodGet, "/test/panic", trt.TestPanic)

	//=========================== User Route
	app.Handle(http.MethodPost, "/user", ugt.Create)
	app.Handle(http.MethodGet, "/user/:id", ugt.QueryByID)
	app.Handle(http.MethodGet, "/user", ugt.Query)
	app.Handle(http.MethodPut, "/user/:id", ugt.Update)
	app.Handle(http.MethodDelete, "/user/:id", ugt.Delete)

	//=========================== User Auth Route
	app.Handle(http.MethodPost, "/user/login", ugt.Login)
	app.Handle(http.MethodPost, "/user/token/refresh", ugt.RefreshToken)
	app.Handle(http.MethodPost, "/user/code/verify", ugt.VerifyConfirmationCode)
	app.Handle(http.MethodGet, "/user/code/resend/:id", ugt.ResendConfirmationCode)
	app.Handle(http.MethodPut, "/user/password/forgot/:id", ugt.ForgotPassword)
	app.Handle(http.MethodPut, "/user/password/reset", ugt.ConfirmNewPassword)

	//=========================== Collaborator Route
	app.Handle(http.MethodPost, "/collaborator", cgt.Create)
	app.Handle(http.MethodGet, "/collaborator/:id", cgt.QueryByID)
	app.Handle(http.MethodGet, "/collaborator", cgt.Query)
	app.Handle(http.MethodPut, "/collaborator/:id", cgt.Update)
	app.Handle(http.MethodDelete, "/collaborator/:id", cgt.Delete)

	//=========================== Collaborator Auth Route
	app.Handle(http.MethodPost, "/collaborator/login", cgt.Login)
	app.Handle(http.MethodPost, "/collaborator/token/refresh", cgt.RefreshToken)
	app.Handle(http.MethodPost, "/collaborator/code/verify", cgt.VerifyConfirmationCode)
	app.Handle(http.MethodGet, "/collaborator/code/resend/:id", cgt.ResendConfirmationCode)
	app.Handle(http.MethodPut, "/collaborator/password/forgot/:id", cgt.ForgotPassword)
	app.Handle(http.MethodPut, "/collaborator/password/reset", cgt.ConfirmNewPassword)

	//=========================== Aggregator Route
	app.Handle(http.MethodPost, "/aggregator", agt.Create)
	app.Handle(http.MethodPut, "/aggregator/:id", agt.Update)
	app.Handle(http.MethodDelete, "/aggregator/:id", agt.Delete)
	app.Handle(http.MethodGet, "/aggregator", agt.Query)
	app.Handle(http.MethodGet, "/aggregator/:id", agt.QueryByID)

	//=========================== Enterprise Route
	app.Handle(http.MethodPost, "/enterprise", egt.Create)
	app.Handle(http.MethodPut, "/enterprise/:id", egt.Update)
	app.Handle(http.MethodDelete, "/enterprise/:id", egt.Delete)
	app.Handle(http.MethodGet, "/enterprise", egt.Query)
	app.Handle(http.MethodGet, "/enterprise/:id", egt.QueryByID)
	app.Handle(http.MethodGet, "/enterprise/code/:code", egt.QueryByCode)

	//=========================== Enterprise Pack Route
	app.Handle(http.MethodPost, "/pack", epgt.Create)
	app.Handle(http.MethodPut, "/pack/:id", epgt.Update)
	app.Handle(http.MethodDelete, "/pack/:id", epgt.Delete)
	app.Handle(http.MethodGet, "/pack", epgt.Query)
	app.Handle(http.MethodGet, "/pack/:id", epgt.QueryByID)

	//=========================== Enterprise Team Route
	app.Handle(http.MethodPost, "/team", tgt.Create)
	app.Handle(http.MethodPut, "/team/:id", tgt.Update)
	app.Handle(http.MethodDelete, "/team/:id", tgt.Delete)
	app.Handle(http.MethodGet, "/team", tgt.Query)
	app.Handle(http.MethodGet, "/team/:id", tgt.QueryByID)
	app.Handle(http.MethodGet, "/team/enterprise/:id", tgt.QueryByEnterprise)

	//=========================== Enterprise Policy Route
	app.Handle(http.MethodPost, "/policy", tpt.Create)
	app.Handle(http.MethodPut, "/policy/:id", tpt.Update)
	app.Handle(http.MethodDelete, "/policy/:id", tpt.Delete)
	app.Handle(http.MethodGet, "/policy", tpt.Query)
	app.Handle(http.MethodGet, "/policy/:id", tpt.QueryByID)
	app.Handle(http.MethodGet, "/policy/enterprise/:id", tpt.QueryByEnterprise)

	//=========================== Invoicing Route
	app.Handle(http.MethodPost, "/invoicing", igt.Create)
	app.Handle(http.MethodPut, "/invoicing/:id", igt.Update)
	app.Handle(http.MethodDelete, "/invoicing/:id", igt.Delete)
	app.Handle(http.MethodGet, "/invoicing", igt.Query)
	app.Handle(http.MethodGet, "/invoicing/:id", igt.QueryByID)
	app.Handle(http.MethodGet, "/invoicing/enterprise/:id", igt.QueryByEnterprise)

	//=========================== Payment Method Route
	app.Handle(http.MethodPost, "/payment/method", pmgt.Create)
	app.Handle(http.MethodPut, "/payment/method/:id", pmgt.Update)
	app.Handle(http.MethodDelete, "/payment/method/:id", pmgt.Delete)
	app.Handle(http.MethodGet, "/payment/method/user/:id", pmgt.Query)

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
