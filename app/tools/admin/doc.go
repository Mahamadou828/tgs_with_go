/*
Admin regroup all cli command to perform admin operations. Below are
the listed admin operations:
- Upload a file in aws s3 bucket
- Migrate database schema
- Upload secrets to aws secret manager
- Seeding the database with mock data

Commands:
	DB:
	- Migrate database
		Params:
			MIGRATE_VERSION - migration version
			ENV 			- aws environment can be: development, staging or production
			AWS_ACCOUNT 	- aws account to use for deployment.
		Command:
			Makefile:
				make db-migrate MIGRATE_VERSION=$(MIGRATE_VERSION)
			Run with go:
				go run app/tools/admin/main.go --commands=db-migrate --version=$(MIGRATE_VERSION) --env=$(ENV) --aws-account=$(AWS_ACCOUNT) | go run app/tools/logfmt/main.go

	- Seed database
			Params:
				MIGRATE_VERSION - migration version
				ENV 			- aws environment can be: development, staging or production
				AWS_ACCOUNT 	- aws account to use for deployment.
			Command:
				go run app/tools/admin/main.go --commands=db-seed --version=$(MIGRATE_VERSION) --env=$(ENV) --aws-account=$(AWS_ACCOUNT) | go run app/tools/logfmt/main.go

	S3:
	- Upload S3 file
		Params:
			VERSION 	- version
			ENV 		- aws environment can be: development, staging or production
			AWS_ACCOUNT - aws account to use for deployment.
			FILE_PATH 	- path to the file to upload
			BUCKET_NAME - name of the bucket where to upload
			BUCKET_KEY 	- the name of the file in s3
		Command:
			MakeFile:
				make s3-upload AWS_ACCOUNT=$(AWS_ACCOUNT) VERSION=$(VERSION) ENV=$(ENV) FILE_PATH=$(FILE_PATH) BUCKET_NAME=$(BUCKET_NAME) BUCKET_KEY=(BUCKET_KEY)
			Run with go:
				go run app/tools/admin/main.go --commands=s3-upload-file --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV) --file=$(FILE_PATH) --bucket=$(BUCKET_NAME) --key=$(BUCKET_KEY)

	Secret Service Manager (SSM):
		- Create secret
				Params:
				ENV 		- aws environment can be: development, staging or production
				AWS_ACCOUNT - aws account to use for deployment.
				SERVICE 	- The service that will use the secret
				SRCTYPE 	- Indicate where the secret will be provided, now can be either: s3, cli, local
				FILENAME 	- optional: the name of the file containing the secrets - must be provided is SrcType is local
				SECRETS 	- optional - list of secrets to create, the format is secrets=name:value:description,name:value:description - must be provided is SrcType is cli
				BUCKET 		- optional - name of the s3 bucket containing the secrets - must be provided is SrcType is s3
				KEY 		- optional - key of the bucket containing the secrets - must be provided is SrcType is s3
			Command:
				go run app/tools/admin/main.go --commands=ssm-create-secrets --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV) --src-type=$(SrcType)

	Api Gateway (agw)
		- Create Spec
			desc: create spec will upload an empty spec to aws s3
			Params:
				ENV 		- aws environment can be: development, staging or production
				AWS_ACCOUNT - aws account to use for deployment.
				SERVICE 	- The service that will use the secret
			Command:
				Makefile:
					make agw-spec-create ENV=$(ENV) AWS_ACCOUNT=$(AWS_ACCOUNT) SERVICE=$(SERVICE)
				Run with go:
					go run app/tools/admin/main.go --commands=agw-spec-create --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV)
		- Create route in spec only
			desc: update the aws spec to add a new route
			Params:
				VERSION 			- version
				ENV 				- aws environment can be: development, staging or production
				AWS_ACCOUNT 		- aws account to use for deployment.
				SERVICE 			- The service that will use the secret
				RESOURCE_NAME		- The agw resource name
				TYPE 				- The route path type, can be: POST, GET, DELETE, PUT
				PATH 				- The route path
				ENABLED_AUTHORIZER 	- is set to true the route will be protected with the lambda authorizer
			Command:
				Makefile:
					make agw-spec-route-create ENV=$(ENV) AWS_ACCOUNT=$(AWS_ACCOUNT) SERVICE=$(SERVICE) TYPE=$(TYPE) RESOURCE_NAME=$(RESOURCE_NAME) PATH=$(PATH) ENABLED_AUTHORIZER=$(ENABLED_AUTHORIZER)
				Run with go:
					go run app/tools/admin/main.go --commands=agw-spec-create --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV)

*/
package main
