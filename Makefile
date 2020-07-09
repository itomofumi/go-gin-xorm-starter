
# make sure BIN_NAME value length is <= 13 charactors.
BIN_NAME=starter-local
BIN_NAME_PRODUCTION=bin/starter

BUILD_TAGS_PRODUCTION='production'
BUILD_TAGS_DEVELOPMENT='development'

ifeq ($(VERSION),)
VERSION := $(GIT_COMMIT_HASH)
endif

# define ANSI color
BOLD=\033[1m
RED=\033[31m
GREEN=\033[32m
CYAN=\033[36m
RESET=\033[0m

# project definition
STORE_ID=1
TZ_OFFSET=-0800
ifeq ($(STORE_ID),2)
TZ_OFFSET := +0900
endif

clean: docker-down
	rm -f debug $(BIN_NAME) $(BIN_NAME_TEST) coverage.out

prepare:
	@if [ ! -e .env ]; then echo "$(RED).env FILE NOT FOUND$(RESET)\n $(GREEN)'cp .env.example .env'$(RESET) to setup .env file." && exit 1; fi

# スタートグループ
watch:
	nodemon -x "pkill $(BIN_NAME) & ($(BUILD_CMD) || exit 1) && (./$(BIN_NAME) || exit 1)"

start: prepare docker-up
	$(MAKE) watch BUILD_CMD="make build-local"

# CI
lint:
	go list ./... | xargs golint -set_exit_status
	go vet ./...

test:
	CGO_ENABLED=0 go test -cover -coverprofile=./coverage.out ./...
	go tool cover -func=coverage.out | grep "total:"

# Docker 操作
docker-up:
	@if [ ! $(docker-compose images | grep aws-cli) ] ; then docker-compose build --parallel ; fi
	@if [ ! $(docker-compose ps | grep mysql) ] ; then docker-compose up -d ; fi

docker-down:
	export CONTAINER_ID=$(shell docker container ls -q -f name=starter-unit-test-mysql14336) ; \
	if [ -n "$${CONTAINER_ID}" ]; then docker container rm -f $${CONTAINER_ID} ; fi
	docker-compose down --remove-orphans

# ビルドグループ

pre-build:
	GO111MODULE=on go build -o $(BIN_NAME) -tags '$(BUILD_TAGS) netgo' -installsuffix netgo \
		-ldflags '-s -w -X github.com/itomofumi/go-gin-xorm-starter/server.version=$(VERSION)' \
		main.go

debug-build:
	GO111MODULE=on go build -o $(BIN_NAME) -tags '$(BUILD_TAGS) netgo' -installsuffix netgo -gcflags 'all=-N -l' \
		-ldflags '-X github.com/itomofumi/go-gin-xorm-starter/server.version=$(VERSION)' \
		main.go

build-local:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_DEVELOPMENT)

