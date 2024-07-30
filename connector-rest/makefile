WORKSPACE ?= $(shell pwd)
include ${WORKSPACE}/.env
export $(shell sed 's/=.*//' .env)
export WORKSPACE
server-build:
	go build -o ${WORKSPACE}/bin/server ${WORKSPACE}/cmd/main.go
server-run:
	${WORKSPACE}/bin/server
build-and-run: server-build server-run
	
