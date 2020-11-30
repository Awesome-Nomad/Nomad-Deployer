SHELL=bash

GIT_IMPORT=github.com/Awesome-Nomad/Nomad-Deployer/cmd
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GOBUILD_VERSION?=$(GIT_COMMIT)$(GIT_DIRTY)
GOLDFLAGS=-s -w -X $(GIT_IMPORT).Version=$(GOBUILD_VERSION)


.PHONY: install
install:
	@go install -ldflags="${GOLDFLAGS}" ./...
	@upx $(shell go env GOPATH)/bin/deployer
.PHONY: build
build:
	go build -ldflags="${GOLDFLAGS}"

.PHONY: test
test:
	@go test -count=1

.PHONY: clean
clean:
	@go clean -testcache -cache

.PHONY: coverage
coverage:
	@go test -short -coverprofile=coverage.out `go list ./... | grep -v vendor/`

