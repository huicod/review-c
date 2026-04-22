GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=internal/conf/conf.proto

# review-c does not maintain local API proto contracts.
# Shared contracts live in ../reviewapis; this service only generates internal/conf.

.PHONY: init
# install proto tooling
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal/conf/conf.proto -> conf.pb.go
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate shared reviewapis proto
api:
	$(MAKE) -C ../reviewapis all

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make config
	make generate

.PHONY: help
# show help
help:
	@echo 'Usage: make [init|config|api|build|generate|all]'

.DEFAULT_GOAL := help