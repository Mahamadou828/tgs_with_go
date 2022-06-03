package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-sdk-go/aws"
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

	//================================================================= VPC
	//create a vpc
	vpc := ec2.NewVpc(scope, &id, &ec2.VpcProps{VpcName: aws.String(fmt.Sprintf("%s-tgs", props.env))})

	//=========================================== Private subnet config

	//create a nat gateway
	
	//create a route table that associat the 0.0.0.0/0 traffic to the nat gateway

	//create a private subnet
	//associate the nat gateway with the private subnet

	//============================================ Public subnet

	//create an internet gateway
	//create a route table that associate the 0.0.0.0/0 with the internet gateway
	//create a public subnet
	//associate the public subnet with the internet gateway

	//================================================================= ECS
	//Get the ecr repository
	//@todo find a way to create it if does not exist and push the first image

	//Create a cluster

	//create a task definition

	//Tell the ECS task to pull Docker image from previously created ECR

	//declare scaling capability

	//create a network load balancer for the cluster

	//create a security group for the cluster allowing only the nlb to send traffic

	//create a database inside the private subnet
	//the database should be accessible only to the cluster app

	//creating a cache

	//create a cognito pool

	//create the sqs queue

	//create bucket invoice

	//create the billing component

	//create the lambda pre-authorizer

	//create api gateway

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
	//NewStack(app, "TgsProductionStack", &TgsStackProps{
	//	awscdk.StackProps{
	//		Env: env(),
	//	},
	//	"production",
	//})
	//NewStack(app, "TgsStagingStack", &TgsStackProps{
	//	awscdk.StackProps{
	//		Env: env(),
	//	},
	//	"staging",
	//})

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
