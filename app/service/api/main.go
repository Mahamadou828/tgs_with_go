package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

//The build represent the environment that the current program is running
//for this specific program we have 3 stages: dev, staging, prod
var build = "dev"

func main() {
	log, err := logger.New("TGS_API")

	if err != nil {
		fmt.Println("can't construct logger")
		panic(err)
	}

	defer log.Sync()

	if err := run(log); err != nil {
		panic(err)
	}
}

func run(log *zap.SugaredLogger) error {
	//===========================
	//GOMAXPROCS

	//Set the correct number of threads for the services
	//based on what is available either by the machine or quotas
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("cpu configuration %w", err)
	}

	log.Info("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	//===========================
	//Init a new aws session
	sesAws, err := aws.New(log)

	if err != nil {
		return fmt.Errorf("can't init an aws session: %w", err)
	}

	//===========================
	//Configuration
	cfg := struct {
		config.Version
		Web struct {
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ApiHost         string        `conf:"default:0.0.0.0:3000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			CorsOrigin      string        `conf:"default:*"`
		}
	}{
		Version: config.Version{
			Build: build,
			Desc:  "TGS api",
		},
	}

	const prefix = "TGS_API"
	help, err := config.Parse(&cfg, prefix, nil)

	if err != nil {
		if errors.Is(err, config.ErrHelpWanted) {
			fmt.Println(help)
		}
		return err
	}

	//===========================
	//App Starting
	log.Infow("starting service", "version", build)
	log.Infow("configuration env", "config", cfg)
	defer log.Infow("shutting down service", "shutting down service", prefix)

	expvar.NewString("build").Set(build)
	expvar.NewString("service").Set(prefix)

	//==========================================================================
	//Start The Debug Server
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	debugMux := handlers.DebugMux(build, log)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown debug router", "status", "debug router error", "host", cfg.Web.DebugHost, "error", err)
		}
	}()

	//==========================================================================
	//Start The Api Server
	log.Infow("initializing", "initializing", "api service starting", "host", cfg.Web.ApiHost)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(web.AppConfig{
		Shutdown:   shutdown,
		Log:        log,
		Build:      build,
		AWS:        sesAws,
		Version:    build,
		Service:    "api",
		CorsOrigin: cfg.Web.CorsOrigin,
	})

	// Construct a server to service the requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.ApiHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}
	//Make a channel to listen for errors coming from the listener. Use a
	//buffered channel so the goroutine can exit if we don't collect this error
	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", cfg.Web.ApiHost)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %v", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
