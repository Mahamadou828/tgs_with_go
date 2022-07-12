//Admin allow to run administrative operations.
package main

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/admin/commands"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"go.uber.org/zap"
)

const service = "admin-cli"

var (
	availableCommands = []string{
		"1. SSM - Create Secret",
		"2. SSM - Create Pool",
		"3. DB - Migrate",
		"4. DB - Seed",
		"5. S3 - Upload JSON file",
	}
)

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
	flag := make(map[string]string)
	if err := config.ParseOsArgs(flag); err != nil {
		return fmt.Errorf("can't parse os args: %v", err)
	}

	cmd, ok := flag["command"]
	if !ok {
		fmt.Println("which command would you like to run?")
		for _, cmd := range availableCommands {
			fmt.Println(cmd)
		}
		fmt.Printf("choice: ")
		if _, err := fmt.Scan(&cmd); err != nil {
			return err
		}
	}

	cfg := struct {
		Service string `conf:"required,help: which service are you targeting?"`
		Env     string `conf:"required,help: which environment are you targeting? (production staging local)"`
		Version string `conf:"required,help: which version are you targeting?"`
	}{}
	if _, err := config.Parse(&cfg, service, parser); err != nil {
		return fmt.Errorf("can't parse %v", err)
	}

	client, err := aws.New(log, aws.Config{
		Service: cfg.Service,
		Env:     cfg.Env,
	})
	if err != nil {
		return fmt.Errorf("can't create aws client: %v", err)
	}

	switch cmd {
	case "1":
		err = commands.SSMCreateSecret(client)
	case "2":
		err = commands.SSMCreatePool(client)
	case "3":
		err = commands.Migrate(
			database.Config{
				User:         "postgres",
				Password:     "postgres",
				Host:         "0.0.0.0:5432",
				Name:         "postgres",
				MaxIdleConns: 0,
				MaxOpenConns: 0,
				DisableTLS:   true,
			},
			cfg.Version,
			log,
		)
	case "4":
		err = commands.Seed(
			database.Config{
				User:         "postgres",
				Password:     "postgres",
				Host:         "0.0.0.0:5432",
				Name:         "postgres",
				MaxIdleConns: 0,
				MaxOpenConns: 0,
				DisableTLS:   true,
			},
			cfg.Version,
			log,
		)
	case "5":
		err = commands.S3UploadJSONFile(client)
	default:
		err = fmt.Errorf("unknown command")
	}
	return err
}

//custom parser for config package
func parser(f config.Field, defaultVal string) error {
	if len(defaultVal) > 0 {
		return config.SetFieldValue(f, defaultVal)
	}
	var val string
	fmt.Printf("%s: ", f.Options.Help)
	if _, err := fmt.Scan(&val); err != nil {
		return fmt.Errorf("invalid value %v", err)
	}
	return config.SetFieldValue(f, val)
}
