SHELL := /bin/bash
KIND_CLUSTER := tgs-cluster

#Vendor all the project dependencies.
tidy:
	go mod tidy
	go mod vendor

#Run the tgs api as a simple go application. Usefull for debugging in local
run-api:
	go run app/service/api/main.go | go run app/tools/logfmt/main.go

VERSION := 1.0

#Start the kind cluster
kind-start: tgs-api kind-up kind-load kind-apply

#Build the docker image for the tgs api
tgs-api:
	docker build \
		-f config/docker/tgs-api.dockerfile \
		-t tgs_api_amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

kind-load:
	cd config/k8s/kind/tgs-pod; kustomize edit set image tgs_api_image=tgs_api_amd64:$(VERSION)
	kind load docker-image tgs_api_amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build config/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build config/k8s/kind/tgs-pod | kubectl apply -f -

kind-up:
	kind create cluster \
		--image kindest/node:v1.23.0@sha256:49824ab1727c04e56a21a5d8372a402fcd32ea51ac96a2706a12af38934f81ac \
		--name $(KIND_CLUSTER) \
		--config config/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=tgs-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubetcl get nodes -o wide
	kubetcl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-logs:
	kubectl logs -l app=tgs --all-containers=true -f --tail=100 | go run app/tools/logfmt/main.go

kind-restart:
	kubetcl rollout restart deployment tgs-pod

kind-update: tgs-api kind-load kind-restart

kind-describe:
	kubectl describe pod -l app=tgs

kind-update-apply: tgs-api kind-load kind-apply

##Use the expvarmon at https://github.com/divan/expvarmon to have local monitoring enabled
##To install expvarmon run: go install github.com/divan/expvarmon@latest
expvarmon:
	expvarmon -ports="4000" -vars="build,requests,gorountines,errors,panics,mem:memstats.Alloc"