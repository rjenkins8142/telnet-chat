## simple makefile
.PHONY: all test build

PROJECT="github.com/rjenkins8142/telnet-chat"
GIT_BUILD=$(shell git rev-parse HEAD)
NOW=$(shell date +%Y-%m-%dT%H:%M:%S)

all: build test

build:
	@go build -ldflags "-X '$(PROJECT)/version.GitBuild=$(GIT_BUILD)' -X '$(PROJECT)/version.BuildDate=$(NOW)'"

test:
	@go test
