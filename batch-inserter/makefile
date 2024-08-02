WORKSPACE ?= $(shell pwd)
include ${WORKSPACE}/.env
export $(shell sed 's/=.*//' .env)
export WORKSPACE
build:
	go build -o ${WORKSPACE}/bin/server ${WORKSPACE}/cmd/main.go
run:
	${WORKSPACE}/bin/server
build-and-run: build run
	
