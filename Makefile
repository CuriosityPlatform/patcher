export APP_SERVICE_CMD_NAME = patcherservice
export APP_CLI_CMD_NAME = patchercli
export REGISTRY = vadimmakerov/curiosity
export APP_PROTO_FILES = \
	api/patcher/patcher.proto
export DOCKER_IMAGE_NAME = $(REGISTRY)-$(APP_SERVICE_CMD_NAME):master

all: build test check

.PHONY: build
build: modules
	bin/go-build.sh "cmd/$(APP_SERVICE_CMD_NAME)" "bin/$(APP_SERVICE_CMD_NAME)" $(APP_SERVICE_CMD_NAME) .env
	bin/go-build.sh "cmd/$(APP_CLI_CMD_NAME)" "bin/$(APP_CLI_CMD_NAME)" $(APP_CLI_CMD_NAME) .env

.PHONY: generate
generate:
	bin/generate-grpc.sh $(foreach path,$(APP_PROTO_FILES),"$(path)")

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run

.PHONY: publish
publish:
	docker build -f data/docker/Dockerfile . --tag=$(DOCKER_IMAGE_NAME)