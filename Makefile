MAKEFLAGS += --silent

.PHONY: build install example fmtcheck fmt testacc

VERSION=0.0.1
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
REGISTRY=registry.terraform.io
NAMESPACE=rode
NAME=rode

build:
	go build -o terraform-provider-rode

install: build
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/
	mv terraform-provider-rode ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}/

example: build
	terraform -chdir=example apply

fmtcheck:
	terraform fmt -recursive -check
	exit $(shell gofmt -s -l . | wc -l)

fmt:
	gofmt -s -w .
	terraform fmt -recursive

test: fmtcheck
	go vet ./...
	go test -v ./...

testacc: build
	TF_ACC=1 RODE_HOST=localhost:50051 RODE_DISABLE_TRANSPORT_SECURITY=true go test ./...