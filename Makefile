GOOS        ?= $(shell go env GOOS)
GOARCH      ?= $(shell go env GOARCH)
GOPATH      ?= $(shell go env GOPATH)
CGO_ENABLED ?= 0

GOBIN := $(if $(shell go env GOBIN),$(shell go env GOBIN),$(GOPATH)/bin)
PATH := $(GOBIN):$(PATH)

COLOR := "\e[1;36m%s\e[0m\n"

define NEWLINE

endef


PROTO_OUT := internal

proto: clean-proto protoc

clean-proto:
	@printf $(COLOR) "Clean old proto files..."
	rm -rf $(PROTO_OUT)/*

protoc: $(PROTO_OUT)
	@printf $(COLOR) "Build proto files..."
	protoc api/v1/service.proto --go_out=plugins=grpc,paths=source_relative:$(PROTO_OUT) api/v1/*.proto