#!/usr/bin/make -f
BINARY = valsta
VERSION ?= $(shell echo $(shell git describe --tags `git rev-list --tags="v*" --max-count=1`) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
build:
ifeq ($(OS),Windows_NT)
	go build -o build/valsta.exe .
else
	go build  -o build/valsta .
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go build -ldflags '$(ldflags)'  -o valsta && mv valsta $(GOPATH)/bin

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs goimports -w -local github.com/jiuyunyue/valsta

setup: build-linux
	@docker build -ldflags '$(ldflags)'  -t valsta .
	@rm -rf ./build


