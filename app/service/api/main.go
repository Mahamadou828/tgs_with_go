/**
@todo add the userPoolID and ClientID inside aws ssm and create a parser to use it
*/
package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Mahamadou828/tgs_with_golang/app/service/api/handlers"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"github.com/Mahamadou828/tgs_with_golang/foundation/web"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

//The build represent the current version of the api
var build = "1.0"

//The env represent the environment that the current program is running
//for this specific program we have 3 stages: dev, staging, prod
var env = "development"

const service = "TGS_API"

type Parser struct {
	Secrets map[string]string
}

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
	sesAws, err := aws.New(log, aws.Config{
		Account: "formation",
		Service: service,
		Env:     env,
	})
	//
	if err != nil {
		return fmt.Errorf("can't init an aws session: %w", err)
	}

	log.Infow("startup", "status", "parsing config struct", "env", env)

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
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres"`
			Host         string `conf:"default:0.0.0.0:5432"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
		Version: config.Version{
			Build: build,
			Desc:  "TGS api",
			Env:   env,
		},
	}

	if env == "staging" || env == "production" {
		secrets, err := sesAws.Ssm.ListSecrets(service, env)

		if err != nil {
			return err
		}

		if help, err := config.Parse(&cfg, service, Parser{Secrets: secrets}); err != nil {
			if errors.Is(err, config.ErrHelpWanted) {
				fmt.Println(help)
			}
			return err
		}
	}

	if env == "development" {
		if help, err := config.Parse(&cfg, service); err != nil {
			if errors.Is(err, config.ErrHelpWanted) {
				fmt.Println(help)
			}
			return err
		}
	}

	//===========================
	//App Starting
	log.Infow("starting service", "version", build)
	log.Infow("configuration env", "config", cfg)
	defer log.Infow("shutting down service", "shutting down service", service)

	expvar.NewString("build").Set(build)
	expvar.NewString("service").Set(service)
	//===========================
	//Open a database connection
	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})

	if err != nil {
		return err
	}

	//==========================================================================
	//Start The Debug Server
	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	debugMux := handlers.DebugMux(build, log, db)

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
		Env:        env,
		Service:    "api",
		CorsOrigin: cfg.Web.CorsOrigin,
		DB:         db,
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
			if err := api.Close(); err != nil {
				return err
			}
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func (p Parser) Parse(field config.Field) error {
	//The value of the field is equal by default to the tag value
	defaultVal := field.Options.DefaultVal

	val, ok := p.Secrets[field.Name]

	//If the secret was not found
	if !ok {
		//And the secret is required we want to terminate the program
		if field.Options.Required {
			return fmt.Errorf("require field %q not present in aws ssm", field.Name)
		}
		//If the secret is not required than we can use the default value
		if !field.Options.Required {
			val = defaultVal
		}
	}

	return config.SetFieldValue(field, val)
}
