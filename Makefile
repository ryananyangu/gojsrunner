VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"
APPLICATION=gojsrunner
PROJECT=integrator


.PHONY: help
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:	## - Build the smallest and secured golang docker image based on scratch (ARGS ARE APPLICATION SPECIFIC)
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Build the smallest and secured golang docker image based on scratch\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build  --build-arg DBURI='$(MongoDbURI)' --build-arg MSGBRD_ACCESSKEY='$(MSGBRD_ACCESSKEY)' -f Dockerfile -t gcr.io/$(PROJECT)/$(APPLICATION):$(VERSION) .

.PHONY: build-no-cache
build-no-cache:	## - Build the smallest and secured golang docker image based on scratch with no cache
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Build the smallest and secured golang docker image based on scratch\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build --no-cache -f Dockerfile -t gcr.io/$(PROJECT)/$(APPLICATION):$(VERSION) .

.PHONY: ls
ls: ## - Listing images of the application with versions
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Listing images of the application with versions !\n\033[0m"
	@docker image ls gcr.io/$(PROJECT)/$(APPLICATION)

.PHONY: run
run:	## - Run the smallest and secured golang docker image based on scratch
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Running image outside environment\n\033[0m"
	@docker run -p 8080:8080 gcr.io/$(PROJECT)/$(APPLICATION):$(VERSION)

.PHONY: run-tests
run-tests: ## - Running tests on environment
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Running tests on environment !\n\033[0m"
	@go test -v -race -coverprofile cover.out  -covermode=atomic ./...
	

.PHONY: run-env
run-env: ## - Running application on local env
	@printf "\033[32m\xE2\x9c\x93 $(VERSION) | Running application on local env !\n\033[0m"
	@docker-compose up -d