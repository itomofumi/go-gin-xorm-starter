
# make sure BIN_NAME value length is <= 13 charactors.
BIN_NAME=starter-local
BIN_NAME_PRODUCTION=bin/starter

BUILD_TAGS_PRODUCTION='production'
BUILD_TAGS_DEVELOPMENT='development'

prepare:
	@if [ ! -e .env ]; then echo "\033[31m.env FILE NOT FOUND\033[0m\n \033[32m'cp .env.example .env'\033[0m to setup .env file." && exit 1; fi

# スタートグループ
start: prepare dep
	nodemon -x "pkill $(BIN_NAME) & (make build-local || exit 1) && (./$(BIN_NAME) || exit 1)"

# CI
lint:
	go list ./... | xargs golint -set_exit_status
	go vet ./...

test:
	go test ./...

# ビルドグループ
dep:
	@dep version
	@dep ensure

pre-build:
	go build -o $(BIN_NAME) -tags '$(BUILD_TAGS) netgo' -installsuffix netgo -ldflags '-s -w' main.go

build-local:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_DEVELOPMENT)

build-linux:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_PRODUCTION) GOOS=linux GOARCH=amd64 BIN_NAME=$(BIN_NAME_PRODUCTION)

build-mac:
	$(MAKE) pre-build BUILD_TAGS=$(BUILD_TAGS_PRODUCTION) GOOS=darwin GOARCH=amd64 BIN_NAME=$(BIN_NAME_PRODUCTION)

