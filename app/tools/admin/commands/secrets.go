package commands

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"go.uber.org/zap"
)

type CreateSecretCfg struct {
	AwsConfig aws.Config
	SrcType   string
	Log       *zap.SugaredLogger
	Filename  string
	Service   string
	Env       string
	Bucket    string
	Key       string
	Secrets   []string
}

type Secret struct {
	Name        string `json:"Name"`
	Value       string `json:"Value"`
	Description string `json:"Description"`
}

//SSMCreateSecret create new secret inside aws secret manager service.
//It's possible to specify three secret src: s3 bucket file, local
//file and cli param. The format of file should be json.
//And the format of cli param should be --secrets=["secretname:secretvalue:secretdesc"]
func SSMCreateSecret(cl *aws.Client) error {
	fmt.Println("Starting creating secret")

	for {
		var name, val string
		fmt.Printf("enter the secret name: ")
		if _, err := fmt.Scan(&name); err != nil {
			return fmt.Errorf("invalid secret name: %v", err)
		}
		fmt.Printf("enter the secret value: ")
		if _, err := fmt.Scan(&val); err != nil {
			return fmt.Errorf("invalid secret value: %v", err)
		}
		fmt.Println("creating new secret")
		if err := cl.SSM.CreateSecret(name, val); err != nil {
			return fmt.Errorf("can't create secret: %v", err)
		}

		fmt.Println("secret created")
		var choice string
		fmt.Printf("Would you like to continue (y|n): ")
		if _, err := fmt.Scan(&choice); err != nil {
			return fmt.Errorf("failed to continue: %v", err)
		}
		if choice == "n" {
			break
		}

	}
	return nil
}

func SSMCreatePool(cl *aws.Client) error {
	if err := cl.SSM.CreatePool(); err != nil {
		return fmt.Errorf("can't create pool: %v", err)
	}
	return nil
}
