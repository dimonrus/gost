TAG=$(shell git branch | sed -n -e 's/\* //p')
PROJECT = $(shell basename `pwd`)
REGISTRY = gost.com/golang

# Clean containers
docker-prune:
	docker system prune -a

# Download golang image
docker-go:
	docker build -t $(REGISTRY)/$(PROJECT)/golang -f .docker/golang.Dockerfile .

# Build main image
docker-build:
	make docker-go
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/build:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/golang -f .docker/build.Dockerfile .

# Build API image
docker-api:
	make docker-swagger
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/api:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/swagger -f .docker/api.Dockerfile .

# Build swagger image
docker-swagger:
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/swagger --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/swagger.Dockerfile .

# Run api image
docker-api-run:
	docker run --rm -p 8080:8080 --env-file="etc/env/.local" --name api $(REGISTRY)/$(PROJECT)/api:$(TAG)

# Build cron image
docker-cron:
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/cron:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/cron.Dockerfile .

# Run cron image
docker-cron-run:
	docker run --rm --env-file="etc/env/.local" --name cron $(REGISTRY)/$(PROJECT)/cron:$(TAG)

# Build test image
docker-test:
	docker build -t $(REGISTRY)/$(PROJECT)/test:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/test.Dockerfile .

# Update test image
docker-test-update:
	docker build -t $(REGISTRY)/$(PROJECT)/test:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/test:$(TAG) -f .docker/test.update.Dockerfile .

# Run tests in docker
docker-test-run:
	docker run --rm --env-file="etc/env/.local" $(REGISTRY)/$(PROJECT)/test:$(TAG)

# Run migration in docker
docker-migrate:
	docker run --rm $(REGISTRY)/$(PROJECT)/build:$(TAG) env ENV=$(TAG) app/$(PROJECT) -app=script -name=migration -class=schema
	docker run --rm $(REGISTRY)/$(PROJECT)/build:$(TAG) env ENV=$(TAG) app/$(PROJECT) -app=script -name=migration -class=data

# Build consumer image
docker-consumer:
	make docker-build
	docker build --no-cache -t $(REGISTRY)/$(PROJECT)/consumer:$(TAG) --build-arg image=$(REGISTRY)/$(PROJECT)/build:$(TAG) -f .docker/consumer.Dockerfile .

# Run consumer image
docker-consumer-run:
	docker run --rm --env-file="etc/env/.local" --name consumer $(REGISTRY)/$(PROJECT)/consumer:$(TAG)

# Create new project
project:
	@read -p "Enter project name: " name; \
	rsync -av --progress . ../$$name --exclude vendor --exclude .idea --exclude .git; \
	find ../$$name gost -path ./vendor -prune -o -path ./.idea -prune -o -path ./.git -prune -o -print -type f -exec sed -i '' -e "s/gost/$$name/g" {} \;

# GO build project
project-build:
	go build -o $(PROJECT) main.go

# Project size
project-size:
	find ./app/ -name '*.go' | xargs wc -l

# Project test count
project-tests:
	env ENV=local go test ./... -v | grep -c RUN

# Create new consumer
script-consumer:
	@read -p "Enter consumer name: " consumer; \
	read -p "Enter queue name: " queue; \
	read -p "Enter server name: " server; \
	env ENV=local ./$(PROJECT) -app=script -name=consumer -file=$$consumer -queue=$$queue -server=$$server;

# Create custom script
script-create:
	@read -p "Enter console action name: " name; \
	env ENV=local ./$(PROJECT) -app script -name create -file $$name

# Create migration
script-migration:
	@read -p "Migration type: " type; \
	read -p "Enter migration name: " name; \
	env ENV=local ./$(PROJECT) -app script -name migration -class $$type -file $$name

# Create model
script-model:
	@read -p "Enter table name: " name; \
	env ENV=local ./$(PROJECT) -app script -name model -file $$name

# Run migration
script-migrate:
	@read -p "Run migration type: " type; \
	env ENV=local ./$(PROJECT) -app script -name migration -class $$type
	
# Download swagger for mac
swagger-mac:
	curl -o swagger  -L https://github.com/go-swagger/go-swagger/releases/download/v0.25.0/swagger_darwin_amd64 && chmod +x swagger

# Download swagger for linux
swagger-lin:
	curl -o swagger  -L https://github.com/go-swagger/go-swagger/releases/download/v0.25.0/swagger_linux_amd64 && chmod +x swagger

# Generate swagger sec
swagger-spec:
	 ./swagger generate spec -m -o swagger.json