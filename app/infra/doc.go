/*
This is the aws infrastructure for the tgs project.

See the infrastructure diagram for more information services used in this project.

**NOTICE**: Go support is still in Developer Preview. This implies that APIs may
change while we address early feedback from the community. We would love to hear
about your experience through GitHub issues.

**Notice**: To make the template work you will need some initial resources:
- You need to create an ecr repository and pushed and initial image, then you will use the repository arn
to create your task definition
- You need a s3 bucket with a json file defining all the api gateway routes. The format file should be:
{
	"ApiKeys": [
		{
			"Value": "",
			"Description": "",
			"UsagePlan": "",
		}
	],
	"Routes": [
		{
			"ResourceName": "paymentMethod",
			"Methods": [
				{
					"Type": "POST",
					"EnabledAuthorizer": false,
					"Path": "/payment/method"
				}
			]
		}
	]
}

## Useful commands

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

*/
package main
