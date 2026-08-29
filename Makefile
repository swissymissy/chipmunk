# makefile to automate capture build command 
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.CommitSHA=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'

.PHONY: build

build:
	mkdir -p dist 
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/chipmunk.exe ./cmd/server