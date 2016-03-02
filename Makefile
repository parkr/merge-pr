RELEASE=$(shell git rev-parse HEAD)

all: build test

deps:
	go get \
	  github.com/bgentry/go-netrc/netrc \
	  github.com/google/go-github/github \
	  golang.org/x/oauth2 \

testdeps:
	go get \
	  github.com/stretchr/testify/assert \
	  golang.org/x/tools/cmd/cover

build: deps
	go build \
	  -ldflags "-X main.release=$(RELEASE)" \
	  -o merge-pr

test: deps testdeps
	go fmt
	go test
