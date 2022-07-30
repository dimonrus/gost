TAG=$(shell git branch | sed -n -e 's/\* //p')
PROJECT = $(shell basename `pwd`)
REGISTRY = gost.com/golang

help:	       			## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

docker-prune:  			## Clean containers
	docker system prune -a

docker-go:     			## Download golang image
	docker build -t $(REGISTRY)/$(PROJECT)/golang -f .docker/golang.Dockerfile .

docker-build:  			## Build main image
	make docker-go
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/build:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/golang -f .docker/build.Dockerfile .

docker-api:			## Build API image
	make docker-swagger
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/api:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/swagger -f .docker/api.Dockerfile .

docker-swagger: 		## Build swagger image
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/swagger --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/swagger.Dockerfile .

docker-api-run:			## Run api image
	docker run --rm -p 8080:8080 --env-file="etc/env/.local" --name api $(REGISTRY)/$(PROJECT)/api:$(TAG)

docker-cron:			## Build cron image
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/cron:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/cron.Dockerfile .

docker-cron-run:		## Run cron image
	docker run --rm --env-file="etc/env/.local" --name cron $(REGISTRY)/$(PROJECT)/cron:$(TAG)

docker-test: 			## Build test image
	docker build -t $(REGISTRY)/$(PROJECT)/test:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/test.Dockerfile .

docker-test-update: 		## Update test image
	docker build -t $(REGISTRY)/$(PROJECT)/test:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/test:$(TAG) -f .docker/test.update.Dockerfile .

docker-test-run:		## Run tests in docker
	docker run --rm --env-file="etc/env/.local" $(REGISTRY)/$(PROJECT)/test:$(TAG)

docker-migrate: 		## Run migration in docker
	docker run --rm $(REGISTRY)/$(PROJECT)/build:$(TAG) env ENV=$(TAG) app/$(PROJECT) -app=script -name=migration -class=schema
	docker run --rm $(REGISTRY)/$(PROJECT)/build:$(TAG) env ENV=$(TAG) app/$(PROJECT) -app=script -name=migration -class=data

docker-consumer: 		## Build consumer image
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/consumer:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/consumer.Dockerfile .

docker-consumer-run:		## Run consumer image
	docker run --rm --env-file="etc/env/.local" --name consumer $(REGISTRY)/$(PROJECT)/consumer:$(TAG)

project: 			## Create new project
	@read -p "Enter project name: " name; \
	rsync -av --progress . ../$$name --exclude vendor --exclude gost --exclude swagger --exclude .idea --exclude .git; \
	find ../$$name gost -path ./vendor -prune -o -path ./.idea -prune -o -path ./.git -prune -o -print -type f -exec sed -i '' -e "s/gost/$$name/g" {} \;

project-build: 			## GO build project
	go build -o $(PROJECT) main.go

project-size: 			## Project size
	find ./app/ -name '*.go' | xargs wc -l

project-tests: 			## Project test count
	env ENV=local go test ./... -v | grep -c RUN

script-consumer: 		## Create new consumer
	@read -p "Enter consumer name: " consumer; \
	read -p "Enter queue name: " queue; \
	read -p "Enter server name: " server; \
	env ENV=local ./$(PROJECT) -app=script -name=consumer -file=$$consumer -queue=$$queue -server=$$server;

script-create: 			## Create custom script
	@read -p "Enter console action name: " name; \
	env ENV=local ./$(PROJECT) -app script -name create -file $$name

script-migration: 		## Create migration
	@read -p "Migration type: " type; \
	read -p "Enter migration name: " name; \
	env ENV=local ./$(PROJECT) -app script -name migration -class $$type -file $$name

script-model: 			## Create model
	@read -p "Enter table name: " name; \
	env ENV=local ./$(PROJECT) -app script -name model -file $$name

script-crud: 			## Create crud
	@read -p "Enter table name: " name; \
	read -p "Enter crud number: " number; \
	env ENV=local ./$(PROJECT) -app script -name crud -file $$name -num $$number

script-migrate:			## Run migration
	@read -p "Run migration type: " type; \
	env ENV=local ./$(PROJECT) -app script -name migration -class $$type

swagger-mac: 			## Download swagger for mac
	curl -o swagger  -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_darwin_amd64 && chmod +x swagger

swagger-lin: 			## Download swagger for linux
	curl -o swagger  -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64 && chmod +x swagger

swagger-spec: 			## Generate swagger sec
	 ./swagger generate spec -m -o swagger.json

tls-pair: 			## Generate tls pair for web server. Required OpenSSL 3.0.0
	@read -p "Enter path: " path; \
	openssl req -new -newkey rsa:4096 -x509 -sha256 -days 1825 -nodes -out $$path/cert.crt -keyout $$path/key.key -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost"

keys:				## Generate rsa pair
	@read -p "Enter keys path: " path; \
	ssh-keygen -t rsa -m PEM -f $$path/id_rsa && openssl rsa -in $$path/id_rsa -pubout -out $$path/id_rsa.pub