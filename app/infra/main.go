package main

import (
	"encoding/json"
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
	"io/ioutil"
	"os"
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
	sess *aws.AWS
}

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

type UsagePlan struct {
	Name        string `json:"name"`
	RateLimit   int    `json:"rateLimit"`
	BurstLimit  int    `json:"burstLimit"`
	Description string `json:"description"`
}

//ApiSpecification regroup the set of route and api key to create along with the api gateway
//To see file format read doc.go file
type ApiSpecification struct {
	ApiKeys    []ApiKey    `json:"apiKeys"`
	Routes     []ApiRoute  `json:"routes"`
	UsagePlans []UsagePlan `json:"usagePlans"`
}

func formatCfn(props *TgsStackProps, name string) string {
	return fmt.Sprintf("%s-%s-%s", props.env, props.service, name)
}

func NewStack(scope constructs.Construct, id *string, props *TgsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, id, &sprops)

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
			SecretName: jsii.String(props.env + "-dbpassword"),
		}),
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
	client := c.AddClient(jsii.String("tgs-api"), &cognito.UserPoolClientOptions{
		AuthFlows: &cognito.AuthFlow{
			AdminUserPassword: jsii.Bool(true),
			Custom:            jsii.Bool(true),
			UserPassword:      jsii.Bool(true),
			UserSrp:           jsii.Bool(true),
		},
		GenerateSecret: jsii.Bool(false),
	})

	awscdk.NewCfnOutput(stack, jsii.String("cognitouserpoolid"), &awscdk.CfnOutputProps{
		Value:       c.UserPoolId(),
		Description: jsii.String(formatCfn(props, "cognitouserpoolid")),
	})

	awscdk.NewCfnOutput(stack, jsii.String("cognitoclientid"), &awscdk.CfnOutputProps{
		Value:       client.UserPoolClientId(),
		Description: jsii.String(formatCfn(props, "cognitoclientid")),
	})

	//generate a seed for the sign-in user preferred_username
	//@todo generate a correct seed
	seed := uuid.NewString()

	awscdk.NewCfnOutput(stack, jsii.String("cognitoseed"), &awscdk.CfnOutputProps{
		Value:       jsii.String(seed),
		Description: jsii.String(formatCfn(props, "cognitoseed")),
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

	//Set the name of the queue inside aws ssm
	if err := props.sess.Ssm.UpdateOrCreateSecret("sqsqueuename", queueName, props.service, props.env, "the name of the sqs queue"); err != nil {
		panic(err)
	}

	//================================================================= S3
	//create bucket invoice
	bucket := s3.NewBucket(stack, jsii.String(props.env+"-tgs-invoices"), &s3.BucketProps{
		Versioned: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(stack, jsii.String("s3invoicesbucketname"), &awscdk.CfnOutputProps{
		Value:       bucket.BucketName(),
		Description: jsii.String(formatCfn(props, "s3invoicesbucketname")),
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

	//Download and parsing the api gateway spec from s3
	file, err := os.Create("spec.json")
	defer file.Close()
	defer os.Remove("spec.json")

	_, err = props.sess.S3.Download(file, "tgs-api-gateway-spec", fmt.Sprintf("spec.%s", props.env))
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var spec ApiSpecification
	if err := json.Unmarshal(bytes, &spec); err != nil {
		panic(err)
	}

	//create the usage plans for the agw
	var usagePlans map[string]agw.UsagePlan

	for _, up := range spec.UsagePlans {
		usagePlans[up.Name] = api.AddUsagePlan(jsii.String(up.Name), &agw.UsagePlanProps{
			ApiStages:   nil,
			Description: jsii.String(up.Description),
			Name:        jsii.String(up.Name),
			//@todo add later
			//Quota:       nil,
			//Throttle:    nil,
		})
	}

	//create the api keys for the agw
	for _, ak := range spec.ApiKeys {
		v, ok := usagePlans[ak.UsagePlan]
		if !ok {
			panic(fmt.Errorf("specified usage plan for api key didn't exist"))
		}

		v.AddApiKey(
			api.AddApiKey(jsii.String(ak.Value), &agw.ApiKeyOptions{
				ApiKeyName:  jsii.String(ak.Name),
				Description: jsii.String(ak.Description),
				Value:       jsii.String(ak.Value),
			}),
			nil,
		)
	}

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

	sess, err := aws.New(log, aws.Config{
		Account:             "formation",
		Service:             "TGS_CDK",
		Env:                 "local",
		UnsafeIgnoreSecrets: true,
	})
	if err != nil {
		log.Errorf("can't init an aws session")
		panic(err)
	}

	NewStack(app, jsii.String("TgsDevelopmentStack"), &TgsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
		"development",
		"TGS_API",
		"arn:aws:ecr:eu-west-1:685367675161:repository/tgs-api-development",
		struct {
			memoryMiB float64
			cpu       float64
		}{memoryMiB: 1024, cpu: 512},
		sess,
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
