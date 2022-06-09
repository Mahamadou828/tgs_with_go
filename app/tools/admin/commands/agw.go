package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	awsToolkit "github.com/aws/aws-sdk-go/aws"
	"go.uber.org/zap"
)

type ApiKey struct {
	Value       string `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UsagePlan   string `json:"usagePlan"`
	GenerateKey bool   `json:"generateKey"`
}

type Method struct {
	Type              string `json:"type"`
	EnabledAuthorizer bool   `json:"enabledAuthorizer"`
	Path              string `json:"path"`
}

type ApiRoute struct {
	ResourceName string   `json:"resourceName"`
	Methods      []Method `json:"methods"`
}

type UsagePlans struct {
	Name        string `json:"name"`
	RateLimit   int    `json:"rateLimit"`
	BurstLimit  int    `json:"burstLimit"`
	Description string `json:"description"`
}

type ApiSpecification struct {
	ApiKeys    []ApiKey     `json:"apiKeys"`
	Routes     []ApiRoute   `json:"routes"`
	UsagePlans []UsagePlans `json:"usagePlans"`
}

const BucketName = "tgs-api-gateway-spec"

//CreateApiKey creates a new API key to the specified api gateway.
func CreateApiKey() {

}

//DeleteApiKey delete an API key from the specified api gateway and update the agw spec file
func DeleteApiKey() {

}

//CreateApiKeyInSpecOnly creates a new API key on the specified api gateway spec
//The api key will not be created unless you run the aws cdk
func CreateApiKeyInSpecOnly() {

}

//CreateAgwSpec creates a new api gateway specification
func CreateAgwSpec(cfg aws.Config, log *zap.SugaredLogger, env string) error {
	sessAws, err := aws.New(log, cfg)
	if err != nil {
		return err
	}

	as := ApiSpecification{
		ApiKeys:    []ApiKey{},
		Routes:     []ApiRoute{},
		UsagePlans: []UsagePlans{},
	}

	b, err := json.Marshal(as)
	if err != nil {
		return err
	}

	if err := sessAws.S3.Upload(b, BucketName, fmt.Sprintf("spec.%s", env), env, "application/json"); err != nil {
		return err
	}

	return nil
}

//CreateRoute creates a new route on the specified api gateway
func CreateRoute(cfg aws.Config, log *zap.SugaredLogger, env string, r ApiRoute) {

}

//CreateRouteInSpecOnly creates a new route only inside the api gateway specification.
//The route will not be deployed unless you run the aws cdk command
func CreateRouteInSpecOnly(cfg aws.Config, log *zap.SugaredLogger, env string, rn string, md Method) error {
	sessAws, err := aws.New(log, cfg)
	as, err := getAgwSpec(cfg, log, env)
	if err != nil {
		return err
	}
	//@todo handle route collision
	as.Routes = append(as.Routes, ApiRoute{Methods: []Method{md}, ResourceName: rn})

	b, err := json.Marshal(as)
	if err != nil {
		return err
	}

	if err := sessAws.S3.Upload(b, BucketName, fmt.Sprintf("spec.%s", env), env, "application/json"); err != nil {
		return err
	}

	return nil
}

func getAgwSpec(cfg aws.Config, log *zap.SugaredLogger, env string) (ApiSpecification, error) {
	sessAws, err := aws.New(log, cfg)
	if err != nil {
		return ApiSpecification{}, err
	}
	buffer := awsToolkit.NewWriteAtBuffer([]byte{})
	if _, err := sessAws.S3.Download(buffer, BucketName, fmt.Sprintf("spec.%s", env)); err != nil {
		return ApiSpecification{}, nil
	}
	b := buffer.Bytes()

	var as ApiSpecification
	if err := json.Unmarshal(b, &as); err != nil {
		return ApiSpecification{}, err
	}

	return as, nil
}
