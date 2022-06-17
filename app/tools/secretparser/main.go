package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"io/ioutil"
	"os"
)

const service = "SECRET_PARSER"

type CfnOutput struct {
	TgsDevelopmentStack struct {
		Cognitoclientid      string `json:"cognitoclientid"`
		S3Invoicesbucketname string `json:"s3invoicesbucketname"`
		Cognitouserpoolid    string `json:"cognitouserpoolid"`
		Cognitoseed          string `json:"cognitoseed"`
		DbHost               string `json:"dbhost"`
		DbUser               string `json:"dbuser"`
		DisableTLS           string `json:"dbdisabletls"`
		DbName               string `json:"dbname"`
		DbPassword           string `json:"dbpassword"`
	} `json:"TgsDevelopmentStack"`
}

func main() {
	log, err := logger.New(service)
	if err != nil {
		panic(err)
	}

	cfg := struct {
		FilePath string `conf:"required"`
	}{}
	log.Infow("parsing configuration structured")
	if h, err := config.Parse(&cfg, service); err != nil {
		if errors.Is(err, config.ErrHelpWanted) {
			fmt.Println(h)
		}
		panic(err)
	}

	log.Infow("open cfn file")
	f, err := os.Open(cfg.FilePath)
	if err != nil {
		panic(err)
	}

	log.Infow("reading cfn file")
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var cfnout CfnOutput

	log.Infow("Unmarshal cfn file")
	if err := json.Unmarshal(bytes, &cfnout); err != nil {
		panic(err)
	}

	log.Infow("initializing new aws session")
	sess, err := aws.New(log, aws.Config{
		Account:             "formation",
		Service:             service,
		UnsafeIgnoreSecrets: true,
	})
	if err != nil {
		panic(err)
	}

	log.Infow("updating cognitoclientid secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"cognitoclientid",
		cfnout.TgsDevelopmentStack.Cognitoclientid,
		"TGS_API",
		"development",
		"cognitoclientid",
	)
	log.Infow("updating cognitouserpoolid secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"cognitouserpoolid",
		cfnout.TgsDevelopmentStack.Cognitouserpoolid,
		"TGS_API",
		"development",
		"cognitouserpoolid",
	)
	log.Infow("updating s3invoicesbucketname secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"s3invoicesbucketname",
		cfnout.TgsDevelopmentStack.S3Invoicesbucketname,
		"TGS_API",
		"development",
		"s3invoicesbucketname",
	)
	log.Infow("updating cognitoseed secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"cognitoseed",
		cfnout.TgsDevelopmentStack.Cognitoseed,
		"TGS_API",
		"development",
		"cognitoclientid",
	)
	log.Infow("updating dbhost secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"dbhost",
		cfnout.TgsDevelopmentStack.DbHost,
		"TGS_API",
		"development",
		"dbhost",
	)
	log.Infow("updating dbuser secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"dbuser",
		cfnout.TgsDevelopmentStack.DbUser,
		"TGS_API",
		"development",
		"dbuser",
	)

	log.Infow("updating dbdisabletls secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"dbdisabletls",
		cfnout.TgsDevelopmentStack.DisableTLS,
		"TGS_API",
		"development",
		"dbdisabletls",
	)
	log.Infow("updating dbname secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"dbname",
		cfnout.TgsDevelopmentStack.DbName,
		"TGS_API",
		"development",
		"dbname",
	)
	log.Infow("updating dbpassword secrets")
	err = sess.Ssm.UpdateOrCreateSecret(
		"dbpassword",
		cfnout.TgsDevelopmentStack.DbPassword,
		"TGS_API",
		"development",
		"dbpassword",
	)

	if err != nil {
		panic(err)
	}
}
