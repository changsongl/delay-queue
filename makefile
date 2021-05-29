GOPATH=$(shell go env GOPATH)

PROGRAM=delayqueue

BINARY=bin/${PROGRAM}
MAIN_FILE=cmd/delayqueue/main.go

VERSION=$(shell git describe --tags --always --long --dirty)
GIT_ADDR=$(shell git remote -v | head -n 1 | awk '{print $$2}')
BUILD_TIME=$(shell date +%FT%T%z)

REPO=github.com/changsongl/delay-queue/

LDFLAGS=-ldflags "-X ${REPO}vars.BuildProgram=${PROGRAM} -X ${REPO}vars.BuildTime=${BUILD_TIME} -X ${REPO}vars.BuildVersion=${VERSION} -X ${REPO}vars.BuildGitPath=${GIT_ADDR}"

build:
	go fmt ./...
#	golint ./...
#	go vet ./...
	go build -o ${BINARY} ${LDFLAGS} ${MAIN_FILE}


#go build -o bin/delayqueue -ldflags "-X delay-queue/vars.BuildProgram=delayqueue" cmd/delayqueue/main.go
