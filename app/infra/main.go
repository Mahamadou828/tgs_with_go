package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type TgsStackProps struct {
	awscdk.StackProps
	env string
}

func NewStack(scope constructs.Construct, id string, props *TgsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	awss3.NewBucket(stack, jsii.String(fmt.Sprintf("testbucket-%s-%s", props.env, id)), nil)
	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewStack(app, "TgsDevelopmentStack", &TgsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		"development",
	})
	NewStack(app, "TgsProductionStack", &TgsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		"production",
	})
	NewStack(app, "TgsStagingStack", &TgsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		"staging",
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Region: jsii.String("eu-west-1"),
	}

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
