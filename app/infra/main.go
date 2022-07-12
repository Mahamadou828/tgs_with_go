package main

import (
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/aws"
	"github.com/Mahamadou828/tgs_with_golang/foundation/logger"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	agw "github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	scalingCapacity "github.com/aws/aws-cdk-go/awscdk/v2/awsapplicationautoscaling"
	cognito "github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecr "github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	ecsPattern "github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	health "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	targets "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2targets"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	sqs "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TgsStackProps struct {
	awscdk.StackProps
	env           string
	service       string
	repositoryARN string
	capacity      struct {
		memoryMiB float64
		cpu       float64
	}
}

func NewStack(scope constructs.Construct, id *string, log *zap.SugaredLogger, props *TgsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, id, &sprops)

	//================================================================= SSM
	//create a pool
	client, err := aws.New(log, aws.Config{Env: props.env, Service: props.service})
	if err != nil {
		panic(err)
	}
	if err := client.SSM.CreatePool(); err != nil {
		panic(err)
	}

	//================================================================= VPC
	//create a vpc with a public and a private subnet
	vpc := ec2.NewVpc(stack, jsii.String(props.env+"vpc"), &ec2.VpcProps{
		VpcName: jsii.String(fmt.Sprintf("%s-tgs", props.env)),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				Name:       jsii.String(props.env + "private-tgs"),
				SubnetType: ec2.SubnetType_PRIVATE_WITH_NAT,
				CidrMask:   jsii.Number(24),
			},
			{
				Name:       jsii.String(props.env + "public-tgs"),
				SubnetType: ec2.SubnetType_PUBLIC,
				CidrMask:   jsii.Number(24),
			},
		},
	})

	//================================================================= ECS
	//Get the ecr repository
	//@todo find a way to create it if does not exist and push the first image
	rep := ecr.Repository_FromRepositoryArn(stack, jsii.String(props.env+"Repository_FromRepositoryArn"), jsii.String(props.repositoryARN))
	//Create a cluster
	clu := ecs.NewCluster(stack, jsii.String(props.env+"NewCluster"), &ecs.ClusterProps{
		ClusterName: jsii.String(fmt.Sprintf("%s-tgs", props.env)),
		Vpc:         vpc,
	})

	//create role for ecs service
	ia := iam.Role_FromRoleArn(
		stack,
		jsii.String(props.env+"Role_FromRoleArn"),
		jsii.String("arn:aws:iam::685367675161:role/TGS_API_SERVICE"),
		nil,
	)

	//create a task definition
	task := ecs.NewFargateTaskDefinition(stack, jsii.String(props.env+"NewTaskDefinition"), &ecs.FargateTaskDefinitionProps{
		TaskRole:       ia,
		Cpu:            jsii.Number(props.capacity.cpu),
		MemoryLimitMiB: jsii.Number(props.capacity.memoryMiB),
	})

	task.AddContainer(jsii.String(props.env+"-tgs-api"), &ecs.ContainerDefinitionOptions{
		//Tell the ECS task to pull Docker image from previously created ECR
		Image:          ecs.ContainerImage_FromEcrRepository(rep, jsii.String("latest")),
		Cpu:            jsii.Number(props.capacity.cpu),
		MemoryLimitMiB: jsii.Number(props.capacity.memoryMiB),
		PortMappings: &[]*ecs.PortMapping{
			{ContainerPort: jsii.Number(3000), HostPort: jsii.Number(3000)},
			{ContainerPort: jsii.Number(4000), HostPort: jsii.Number(4000)},
		},
	})

	//create new security groups

	//create an application load balancer for the cluster
	alb := ecsPattern.NewApplicationLoadBalancedFargateService(
		stack,
		jsii.String(props.env+"NewApplicationLoadBalancedFargateService"),
		&ecsPattern.ApplicationLoadBalancedFargateServiceProps{
			Cluster:            clu,
			MinHealthyPercent:  jsii.Number(50),
			MaxHealthyPercent:  jsii.Number(300),
			PublicLoadBalancer: jsii.Bool(false),
			ServiceName:        jsii.String(fmt.Sprintf("%s-tgs", props.env)),
			Cpu:                jsii.Number(1024),
			MemoryLimitMiB:     jsii.Number(2048),
			TaskDefinition:     task,
			DesiredCount:       jsii.Number(1),
		})

	//allow traffic to port 4000 for health checks
	alb.LoadBalancer().Connections().AllowFromAnyIpv4(ec2.Port_Tcp(jsii.Number(4000)), jsii.String("allow inbound https"))
	alb.LoadBalancer().Connections().AllowToAnyIpv4(ec2.Port_Tcp(jsii.Number(4000)), jsii.String("allow outbound https"))
	alb.Service().Connections().AllowFromAnyIpv4(ec2.Port_Tcp(jsii.Number(4000)), jsii.String("allow inbound https"))
	alb.Service().Connections().AllowToAnyIpv4(ec2.Port_Tcp(jsii.Number(4000)), jsii.String("allow outbound https"))

	//configure health check
	alb.TargetGroup().ConfigureHealthCheck(&health.HealthCheck{
		Enabled: jsii.Bool(true),
		Path:    jsii.String("/debug/liveliness"),
		Port:    jsii.String("4000"),
	})

	//configure scaling policy
	as := alb.Service().AutoScaleTaskCount(&scalingCapacity.EnableScalingProps{
		MaxCapacity: jsii.Number(10),
		MinCapacity: jsii.Number(1),
	})

	as.ScaleOnMemoryUtilization(
		jsii.String(fmt.Sprintf("%s-ScaleOnMemoryUtilization", props.env)),
		&ecs.MemoryUtilizationScalingProps{
			TargetUtilizationPercent: jsii.Number(70),
		},
	)

	//Create a network load balancer to forward request to alb
	nlb := health.NewNetworkLoadBalancer(stack, jsii.String("Nlb"), &health.NetworkLoadBalancerProps{
		Vpc:            vpc,
		InternetFacing: jsii.Bool(false),
	})

	listener := nlb.AddListener(jsii.String("listener"), &health.BaseNetworkListenerProps{
		Port: jsii.Number(80),
	})

	t := listener.AddTargets(jsii.String("Targets"), &health.AddNetworkTargetsProps{
		Targets: &[]health.INetworkLoadBalancerTarget{
			targets.NewAlbTarget(alb.LoadBalancer(), jsii.Number(80)),
		},
		Port:        jsii.Number(80),
		HealthCheck: &health.HealthCheck{Protocol: health.Protocol_HTTP},
	})

	//see: https://github.com/aws/aws-cdk/issues/17208
	t.Node().AddDependency(alb.Listener())

	////================================================================= Database
	rdsPoolName := fmt.Sprintf("%sRdsPool", props.env)
	////create a database inside the private subnet
	rds.NewDatabaseInstance(stack, jsii.String(props.env+"database"), &rds.DatabaseInstanceProps{
		Vpc:           vpc,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		VpcSubnets: &ec2.SubnetSelection{
			SubnetType: ec2.SubnetType_PUBLIC,
		},
		Engine: rds.DatabaseInstanceEngine_Postgres(&rds.PostgresInstanceEngineProps{
			Version: rds.PostgresEngineVersion_VER_14_2(),
		}),
		DatabaseName: jsii.String(props.env + "database"),
		InstanceType: ec2.InstanceType_Of(ec2.InstanceClass_BURSTABLE3, ec2.InstanceSize_SMALL),
		Credentials: rds.Credentials_FromGeneratedSecret(jsii.String(props.env), &rds.CredentialsBaseOptions{
			SecretName: jsii.String(rdsPoolName),
		}),
	})

	awscdk.NewCfnOutput(stack, jsii.String("rdsPoolName"), &awscdk.CfnOutputProps{
		Value: jsii.String(rdsPoolName),
	})

	//================================================================= Cognito
	//create a cognito pool
	c := cognito.NewUserPool(stack, jsii.String(props.env+"cognitopool"), &cognito.UserPoolProps{
		UserPoolName:  jsii.String(props.env + "-tgs-cognito"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		SignInAliases: &cognito.SignInAliases{
			Username:          jsii.Bool(true),
			PreferredUsername: jsii.Bool(true),
		},
		AutoVerify: &cognito.AutoVerifiedAttrs{
			Phone: jsii.Bool(true),
		},
		StandardAttributes: &cognito.StandardAttributes{
			Email: &cognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(true),
			},
			PhoneNumber: &cognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(true),
			},
			Fullname: &cognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(true),
			},
		},
		PasswordPolicy: &cognito.PasswordPolicy{
			MinLength:        jsii.Number(12),
			RequireLowercase: jsii.Bool(true),
			RequireUppercase: jsii.Bool(true),
			RequireDigits:    jsii.Bool(true),
			RequireSymbols:   jsii.Bool(true),
		},
		CustomAttributes: &map[string]cognito.ICustomAttribute{
			"isActive": cognito.NewStringAttribute(&cognito.StringAttributeProps{
				MinLen:  jsii.Number(1),
				MaxLen:  jsii.Number(256),
				Mutable: jsii.Bool(true),
			}),
			"aggregator": cognito.NewStringAttribute(&cognito.StringAttributeProps{
				MinLen:  jsii.Number(1),
				MaxLen:  jsii.Number(256),
				Mutable: jsii.Bool(false),
			}),
		},
		AccountRecovery: cognito.AccountRecovery_PHONE_ONLY_WITHOUT_MFA,
	})

	//Create a new App client
	poolClient := c.AddClient(jsii.String("tgs-api"), &cognito.UserPoolClientOptions{
		AuthFlows: &cognito.AuthFlow{
			AdminUserPassword: jsii.Bool(true),
			Custom:            jsii.Bool(true),
			UserPassword:      jsii.Bool(true),
			UserSrp:           jsii.Bool(true),
		},
		GenerateSecret: jsii.Bool(false),
	})

	awscdk.NewCfnOutput(stack, jsii.String("cognitoUserPoolId"), &awscdk.CfnOutputProps{
		Value: c.UserPoolId(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("cognitoClientPoolId"), &awscdk.CfnOutputProps{
		Value: poolClient.UserPoolClientId(),
	})

	//generate a seed for the sign-in user preferred_username
	//@todo generate a correct seed
	seed := uuid.NewString()

	awscdk.NewCfnOutput(stack, jsii.String("cognitoSeed"), &awscdk.CfnOutputProps{
		Value: jsii.String(seed),
	})

	//@todo add lambda trigger for confirm signup
	//create the lambda pre-authorizer

	//================================================================= Queue
	//create the sqs queue
	queueName := props.env + "-tgs-queue"
	sqs.NewQueue(stack, jsii.String(props.env+"-tgs-queue"), &sqs.QueueProps{
		DeliveryDelay:     awscdk.Duration_Seconds(jsii.Number(15)),
		QueueName:         jsii.String(queueName),
		VisibilityTimeout: awscdk.Duration_Hours(jsii.Number(12)),
	})

	awscdk.NewCfnOutput(stack, jsii.String("sqsQueueName"), &awscdk.CfnOutputProps{
		Value: jsii.String(queueName),
	})

	//================================================================= S3
	//create bucket invoice
	bucket := s3.NewBucket(stack, jsii.String(props.env+"-tgs-invoices"), &s3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(stack, jsii.String("s3InvoiceBucketName"), &awscdk.CfnOutputProps{
		Value: bucket.BucketName(),
	})

	//================================================================= Billing component
	//create the billing component

	//================================================================= API Gateway
	//create a vpc link
	link := agw.NewVpcLink(stack, jsii.String(props.env+"tgs-vpc-link"), &agw.VpcLinkProps{
		Targets: &[]health.INetworkLoadBalancer{
			nlb,
		},
	})

	//create api gateway
	api := agw.NewRestApi(stack, jsii.String(props.env+"-tgs"), &agw.RestApiProps{
		//Enable Cors
		DefaultCorsPreflightOptions: &agw.CorsOptions{
			AllowOrigins: &[]*string{jsii.String("*")},
			AllowHeaders: &[]*string{
				jsii.String("Content-Type"),
				jsii.String("Authorization"),
				jsii.String("X-Amz-Date"),
				jsii.String("X-Api-Key"),
				jsii.String("x-api-key"),
				jsii.String("X-Amz-Security-Token"),
				jsii.String("aggregatorCode"),
			},
			AllowMethods: &[]*string{
				jsii.String("POST"),
				jsii.String("PUT"),
				jsii.String("GET"),
				jsii.String("OPTIONS"),
				jsii.String("DELETE"),
			},
			StatusCode: jsii.Number(200),
		},
		Deploy: jsii.Bool(true),
		DeployOptions: &agw.StageOptions{
			StageName: jsii.String("main"),
			//If needed you can specify stage variables here
			Variables: nil,
		},
		RestApiName:       jsii.String(props.env + "-tgs"),
		RetainDeployments: jsii.Bool(true),
		Description:       jsii.String("an api gateway to access the tgs api in " + props.env),
	})

	//create a proxy resource
	proxyResource := api.Root().AddProxy(&agw.ProxyResourceOptions{
		AnyMethod: jsii.Bool(false),
	})

	proxyResource.AddMethod(
		jsii.String("ANY"),
		agw.NewIntegration(&agw.IntegrationProps{
			Type: agw.IntegrationType_HTTP_PROXY,
			Options: &agw.IntegrationOptions{
				ConnectionType: agw.ConnectionType_VPC_LINK,
				VpcLink:        link,
			},
			IntegrationHttpMethod: jsii.String("ANY"),
		}),
		&agw.MethodOptions{
			ApiKeyRequired: jsii.Bool(true),
		},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	log, err := logger.New("TGS_CDK")
	if err != nil {
		log.Errorf("failed to create an logger")
		panic(err)
	}

	NewStack(
		app,
		jsii.String("development"),
		log,
		&TgsStackProps{
			awscdk.StackProps{
				Env: env(),
			},
			"development",
			"tgs-api",
			"arn:aws:ecr:eu-west-1:685367675161:repository/tgs-api-development",
			struct {
				memoryMiB float64
				cpu       float64
			}{memoryMiB: 1024, cpu: 512},
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

// env determines the Client environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.jsii.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Region: jsii.String("eu-west-1"),
	}
}
