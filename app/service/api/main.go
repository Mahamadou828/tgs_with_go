package main

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws/session"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws/ssm"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
)

//The build represent the environment that the current program is running
//for this specific programm we have 3 stages: dev, staging, prod
var build = "dev"

func main() {
	log, err := logger.NewLogger("tgs-api")

	if err != nil {
		panic(err)
	}

	log.Info("Testing the logger package")

	//=====================Initiate new AWS Session ================//
	_, err = session.New()

	if err != nil {
		panic(err)
	}

	//=====================Testing ssm package================//
	secret, err := ssm.ListSecrets("tgs-api", build)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Secrets: %v\n", secret)

	//=====================Testing config package================//
	ssmSecrets, err := ssm.ListSecrets("tgs-api", build)
	cfg := struct {
		Web struct {
			port int `conf:"default:3000"`
		}
		DB struct {
			host string `conf:"default:localhost:3000"`
		}
	}{}
	config.Parse(&cfg, ssmSecrets, "TGS_API")
}
