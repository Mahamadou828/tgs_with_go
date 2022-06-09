##@todo when launching the api in other env than dev, the app can't connect to the aws service because we should mount the local .aws file to a volume inside the kluster
SHELL := /bin/bash
KIND_CLUSTER := tgs-cluster
ENV := development
VERSION := 1.0
DB_VERSION := v1
AWS_ACCOUNT := formation

#Vendor all the project dependencies.
tidy:
	go mod tidy
	go mod vendor

#Run the tgs api as a simple go application. Useful for debugging in local
#If you want to start different environments use the following command:
#make ENV="env" run-api
#Be aware that the migration is run automatically so be careful when running another env than development
run-api:
	docker run --name postgres-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres || true
	go run -ldflags "-X main.build=${VERSION}" -ldflags "-X main.env=${ENV}" app/service/api/main.go | go run app/tools/logfmt/main.go



#Start the kind cluster in development mode by default.
#To start the kind cluster in other mode use:
#make ENV="env" kind-start
kind-start: tgs-api kind-up kind-load kind-apply


#Build the docker image for the tgs api
#By default, the api run on development env
#to set the env use
#make ENV='env' tgs-api
tgs-api:
	docker build \
		-f config/docker/tgs-api.dockerfile \
		-t tgs_api_amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg ENV=$(ENV) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

kind-load:
	cd config/k8s/kind/tgs-pod; kustomize edit set image tgs_api_image=tgs_api_amd64:$(VERSION)
	kind load docker-image tgs_api_amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build config/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=360s --for=condition=Available deployment/database-pod
	kustomize build config/k8s/kind/tgs-pod | kubectl apply -f -

kind-up:
	kind create cluster \
		--image kindest/node:v1.23.0@sha256:49824ab1727c04e56a21a5d8372a402fcd32ea51ac96a2706a12af38934f81ac \
		--name $(KIND_CLUSTER) \
		--config config/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=staging

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-logs:
	kubectl logs -l app=tgs --all-containers=true -f --tail=100 | go run app/tools/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment tgs-pod

kind-update: tgs-api kind-load kind-restart

kind-describe:
	kubectl describe pod -l app=tgs

kind-update-apply: tgs-api kind-load kind-apply

##Use the expvarmon at https://github.com/divan/expvarmon to have local monitoring enabled
##To install expvarmon run: go install github.com/divan/expvarmon@latest
expvarmon:
	expvarmon -ports="3000" -vars="build,requests,gorountines,errors,panics,mem:memstats.Alloc"

db-start: db-up db-migrate db-seed

db-restart: db-destroy db-up db-migrate db-seed

#================================================================= Admin Command
#Destroy the local postgres sql database
db-destroy:
	docker stop postgres-db || true && docker rm postgres-db || true

#Create a local postgresql database
db-up:
	docker run --name postgres-db -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

#Migrate the database schemas to use the command you should provide the MIGRATE_VERSION:
#make db-migration DB_VERSION=v1
db-migrate:
	go run app/tools/admin/main.go --commands=migrate --version=$(DB_VERSION) --env=$(ENV) --awsaccount=$(AWS_ACCOUNT) | go run app/tools/logfmt/main.go

#Seed the database with fake data:
#make db-seed DB_VERSION=v1
db-seed:
	go run app/tools/admin/main.go --commands=seed --version=$(DB_VERSION) --env=$(ENV) --awsaccount=$(AWS_ACCOUNT) | go run app/tools/logfmt/main.go
#Upload a s3 file
#make s3-upload ENV=development FILE_PATH=FILE_PATH BUCKET_NAME=BUCKET_NAME BUCKET_KEY=BUCKET_KEY
s3-upload:
	go run app/tools/admin/main.go --commands=s3-upload-file --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV) --file=$(FILE_PATH) --bucket=$(BUCKET_NAME) --key=$(BUCKET_KEY)

#Create Api Gateway spec
#make agw-create-spec ENV=$(ENV) AWS_ACCOUNT=$(AWS_ACCOUNT) VERSION=$(VERSION)
agw-spec-create:
	go run app/tools/admin/main.go --commands=agw-spec-create --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV)


#Create a new agw route inside the spec
#make agw-spec-route-create ENV=$(ENV) AWS_ACCOUNT=$(AWS_ACCOUNT) VERSION=$(VERSION) TYPE=$(TYPE) RESOURCE_NAME=$(RESOURCE_NAME) RESOURCE_PATH=$(RESOURCE_PATH) ENABLED_AUTHORIZER=$(ENABLED_AUTHORIZER)
agw-spec-route-create:
	go run app/tools/admin/main.go --commands=agw-spec-route-create --aws-account=$(AWS_ACCOUNT) --version=$(VERSION) --env=$(ENV) --type=$(TYPE) --resource-name=$(RESOURCE_NAME) --resource-path=$(RESOURCE_PATH) --enabled-authorizer=$(ENABLED_AUTHORIZER)
