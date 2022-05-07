//Package web contains a custom web framework that
//wrap the https://github.com/dimfeld/httptreemux package
package web

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/dimfeld/httptreemux"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//A Handler is a type that handles a Http requests withing our framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) *RequestError

type AppConfig struct {
	Shutdown   chan os.Signal
	Log        *zap.SugaredLogger
	Build      string
	AWS        *aws.AWS
	Env        string
	Service    string
	CorsOrigin string
	DB         *sqlx.DB
}

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	mux      *httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
	AWS      *aws.AWS
	group    *httptreemux.ContextGroup
}

//NewApp creates a new App value that handle a set of routes for the application
func NewApp(cfg AppConfig, version string, mw ...Middleware) *App {
	mux := httptreemux.NewContextMux()

	group := mux.NewGroup(fmt.Sprintf("/%s/%s", version, cfg.Service))

	return &App{
		mux:      mux,
		shutdown: cfg.Shutdown,
		mw:       mw,
		AWS:      cfg.AWS,
		group:    group,
	}
}

func (a *App) Handle(method, path string, handler Handler, mw ...Middleware) {
	//If a set of middleware was pass we should wrap it around the main handler
	if mw != nil {
		handler = wrapMiddleware(mw, handler)
	}

	//Second we wrap the app level middleware
	handler = wrapMiddleware(a.mw, handler)

	//The function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//For now let's use a false traceID
		//@todo generate a unique traceid for sentry logs
		span := uuid.NewString()

		v := RequestTrace{
			ID:  span,
			Now: time.Now().UTC(),
		}

		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	a.group.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
