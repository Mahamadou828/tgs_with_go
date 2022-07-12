package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/app/tools/config"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

const service = "SECRET_PARSER"

func main() {
	log, err := logger.New(service)
	if err != nil {
		fmt.Printf("cannot init logger: %v", err)
		panic(err)
	}
	cfg := struct {
		Service  string `conf:"required"`
		FilePath string `conf:"required"`
	}{}
	if _, err := config.Parse(&cfg, service, nil); err != nil {
		log.Infof("failed to parse configuration %v", err)
		os.Exit(1)
	}

	jsonFile, err := os.Open(cfg.FilePath)
	if err != nil {
		log.Infof("failed to open file %v", err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Infof("failed to read file %v", err)
		os.Exit(1)
	}

	var out map[string]map[string]string
	if err := json.Unmarshal(b, &out); err != nil {
		log.Infof("failed to unmarshal file %v", err)
		os.Exit(1)
	}

	for env, secretMap := range out {
		var secrets []aws.Secret
		log.Infow("creating client for env", "env", env)
		log.Infow("starting to create secret")
		client, err := aws.New(log, aws.Config{
			Service: cfg.Service,
			Env:     env,
		})
		if err != nil {
			log.Infof("failed to create an client for env %s: %v", env, err)
			continue
		}
		for key, val := range secretMap {
			log.Infow("creating secret", "name", key)
			if key == "rdsPoolName" {
				log.Infow("detect special secret name rdsPoolName. Start pooling secret for re-creation...")
				r := struct {
					Password             string `json:"password"`
					DbName               string `json:"dbname"`
					Engine               string `json:"-"`
					Port                 int    `json:"port"`
					DbInstanceIdentifier string `json:"-"`
					Host                 string `json:"host"`
					Username             string `json:"username"`
				}{}
				if err := client.SSM.GetPoolSecrets(val, &r); err != nil {
					log.Infof("can't get pool secrets: %v", err)
					os.Exit(1)
				}
				secrets = append(
					secrets,
					aws.Secret{Name: "RDS_DB_PASSWORD", Value: r.Password},
					aws.Secret{Name: "RDS_DB_NAME", Value: r.DbName},
					aws.Secret{Name: "RDS_DB_PORT", Value: strconv.Itoa(r.Port)},
					aws.Secret{Name: "RDS_DB_HOST", Value: r.Host},
					aws.Secret{Name: "RDS_DB_USER", Value: r.Username},
				)
				log.Infow("rdsPoolName secret pool")
				continue
			}
			secrets = append(secrets, aws.Secret{Name: ToSnake(key), Value: val})

		}
		log.Infow("creating secret inside pool", "poolname", fmt.Sprintf("%s/%s", cfg.Service, env))
		err = client.SSM.CreateSecrets(secrets)
		if err != nil {
			log.Infof("can't create database secrets: %v", err)
			os.Exit(1)
		}
	}
}

func ToSnake(camel string) (snake string) {
	snake = matchAllCap.ReplaceAllString(matchFirstCap.ReplaceAllString(camel, "${1}_${2}"), "${1}_${2}")
	return strings.ToUpper(snake)
}
