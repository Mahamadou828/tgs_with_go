#start the local postgres database
docker-up:
	docker-compose -f ./config/docker/docker-compose.yml up --force-recreate -d
#destroy the local postgres database. WARNING: this will delete all the data store and free the volume
docker-down:
	docker-compose -f ./config/docker/docker-compose.yml down

#Vendor all the project dependencies.
tidy:
	go mod tidy
	go mod vendor

run-api:
	go run app/service/api/main.go | go run app/tools/logfmt/main.go