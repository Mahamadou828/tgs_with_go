//Admin allow to run administrative operations.
package main

import (
	"errors"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/admin/commands"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"go.uber.org/zap"
)

const service = "ADMIN_SERVICE"

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
	//Configuration
	cfg := struct {
		Commands   []string `conf:"required"`
		Version    string   `conf:"required"`
		Env        string   `conf:"default:development"`
		AwsAccount string   `conf:"required"`
	}{}

	help, err := config.Parse(&cfg, service)

	if err != nil {
		if errors.Is(err, config.ErrHelpWanted) {
			fmt.Println(help)
			return err
		}

		return err
	}

	for _, command := range cfg.Commands {
		switch command {
		case "migrate":
			dbCfg := struct {
				User         string `conf:"default:postgres"`
				Password     string `conf:"default:postgres"`
				Host         string `conf:"default:0.0.0.0:5432"`
				Name         string `conf:"default:postgres"`
				MaxIdleConns int    `conf:"default:0"`
				MaxOpenConns int    `conf:"default:0"`
				DisableTLS   bool   `conf:"default:true"`
			}{}

			if _, err = config.Parse(&dbCfg, service); err != nil {
				return fmt.Errorf("can't start command %s because of missing configuration %w", command, err)
			}
			err = commands.Migrate(
				database.Config{
					User:         dbCfg.User,
					Password:     dbCfg.Password,
					Host:         dbCfg.Host,
					Name:         dbCfg.Name,
					MaxIdleConns: dbCfg.MaxIdleConns,
					MaxOpenConns: dbCfg.MaxOpenConns,
					DisableTLS:   dbCfg.DisableTLS,
				},
				cfg.Version,
				log,
			)
		case "seed":
			dbCfg := struct {
				User         string `conf:"default:postgres"`
				Password     string `conf:"default:postgres"`
				Host         string `conf:"default:0.0.0.0:5432"`
				Name         string `conf:"default:postgres"`
				MaxIdleConns int    `conf:"default:0"`
				MaxOpenConns int    `conf:"default:0"`
				DisableTLS   bool   `conf:"default:true"`
			}{}

			if _, err = config.Parse(&dbCfg, service); err != nil {
				return fmt.Errorf("can't start command %s because of missing configuration %w", command, err)
			}
			err = commands.Seed(
				database.Config{
					User:         dbCfg.User,
					Password:     dbCfg.Password,
					Host:         dbCfg.Host,
					Name:         dbCfg.Name,
					MaxIdleConns: dbCfg.MaxIdleConns,
					MaxOpenConns: dbCfg.MaxOpenConns,
					DisableTLS:   dbCfg.DisableTLS,
				},
				cfg.Version,
				log,
			)
		case "createsecret":
			scrCfg := struct {
				Filename string `conf:"optional"`
				Service  string `conf:"required"`
				Bucket   string `conf:"optional"`
				Key      string `conf:"optional"`
				SrcType  string `conf:"optional"`
			}{}

			if _, err = config.Parse(&scrCfg, service); err != nil {
				return fmt.Errorf("can't start command %s because of missing configuration %w", command, err)
			}

			err = commands.CreateSecret(commands.CreateSecretCfg{
				AwsConfig: aws.Config{
					Account:             cfg.AwsAccount,
					Service:             service,
					Env:                 cfg.Env,
					UnsafeIgnoreSecrets: true,
				},
				SrcType:  scrCfg.SrcType,
				Log:      log,
				Filename: scrCfg.Filename,
				Service:  scrCfg.Service,
				Env:      cfg.Env,
				Bucket:   scrCfg.Bucket,
				Key:      scrCfg.Key,
			})
		case "uploadfile":
			uplCfg := struct {
				File   string `conf:"required"`
				Bucket string `conf:"required"`
				Key    string `conf:"required"`
			}{}

			if _, err = config.Parse(&uplCfg, service); err != nil {
				return fmt.Errorf("can't start command %s because of missing configuration %w", command, err)
			}

			err = commands.Download(
				aws.Config{
					Account:             cfg.AwsAccount,
					Service:             service,
					Env:                 cfg.Env,
					UnsafeIgnoreSecrets: true,
				},
				log,
				uplCfg.File,
				uplCfg.Bucket,
				uplCfg.Key,
			)
		case "test":
			fmt.Println("Test Command")
		}
	}

	return err
}
