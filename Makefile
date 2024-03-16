SHELL := /bin/bash
PWD := $(shell pwd)

GIT_REMOTE = github.com/7574-sistemas-distribuidos/docker-compose-init

default: build

all:

deps:
	go mod tidy
	go mod vendor

build: deps
	GOOS=linux go build -o bin/client github.com/7574-sistemas-distribuidos/docker-compose-init/client
.PHONY: build

docker-image:
	docker build -f ./server/Dockerfile -t "server:latest" .
	docker build -f ./client/Dockerfile -t "client:latest" .
	docker build --network host -f ./netcat-client/Dockerfile -t "netcat-client:latest" .
	# Execute this command from time to time to clean up intermediate stages generated 
	# during client build (your hard drive will like this :) ). Don't left uncommented if you 
	# want to avoid rebuilding client image every time the docker-compose-up command 
	# is executed, even when client code has not changed
	# docker rmi `docker images --filter label=intermediateStageToBeDeleted=true -q`
.PHONY: docker-image

docker-compose-up:
	@echo "Generating Docker Compose configuration for $(clients) clients..."
	chmod u+x docker-compose-script-with-n-clients.sh
	./docker-compose-script-with-n-clients.sh $(clients) > docker-compose-n-clients.yaml
	$(MAKE) docker-image
	docker compose -f docker-compose-n-clients.yaml up -d --build --remove-orphans
.PHONY: docker-compose-up

docker-compose-down:
	docker compose -f docker-compose-n-clients.yaml stop -t 10
	docker compose -f docker-compose-n-clients.yaml down
.PHONY: docker-compose-down

docker-compose-logs:
	docker compose -f docker-compose-n-clients.yaml logs -f
.PHONY: docker-compose-logs

run-netcat-test:
	docker build --network host -f ./netcat-client/Dockerfile -t "netcat-client:latest" .
	$(MAKE) docker-compose-up clients=0
	docker run --rm --network sistemas-distribuidos-tp0_testing_net --env-file ./netcat-client/config.env --name netcat-client netcat-client:latest