
# make sure BIN_NAME value length is <= 13 charactors.
BIN_NAME=starter-local
BIN_NAME_PRODUCTION=bin/starter

BUILD_TAGS_PRODUCTION='production'
BUILD_TAGS_DEVELOPMENT='development'

# define ANSI color
BOLD=\033[1m
RED=\033[31m
GREEN=\033[32m
CYAN=\033[36m
RESET=\033[0m

prepare:
	@if [ ! -e .env ]; then echo "$(RED).env FILE NOT FOUND$(RESET)\n $(GREEN)'cp .env.example .env'$(RESET) to setup .env file." && exit 1; fi

# スタートグループ
start: prepare dep
	nodemon -x "pkill $(BIN_NAME) & (make build-local || exit 1) && (./$(BIN_NAME) || exit 1)"

# CI
lint:
	go list ./... | xargs golint -set_exit_status
	go vet ./...

test:
	go test -cover -race ./...

test-coverage:
	CGO_ENABLED=0 go test -cover -coverprofile=./coverage.out ./...

# ビルドグループ
dep:
	@dep version
	@dep ensure

pre-build:
	go build -o $(BIN_NAME) -tags '$(BUILD_TAGS) netgo' -installsuffix netgo -ldflags '-s -w' main.go

build-local:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_DEVELOPMENT)

.PHONY: build
build:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_PRODUCTION) BIN_NAME=$(BIN_NAME_PRODUCTION)

build-linux:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_PRODUCTION) GOOS=linux GOARCH=amd64 BIN_NAME=$(BIN_NAME_PRODUCTION)

build-mac:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_PRODUCTION) GOOS=darwin GOARCH=amd64 BIN_NAME=$(BIN_NAME_PRODUCTION)

