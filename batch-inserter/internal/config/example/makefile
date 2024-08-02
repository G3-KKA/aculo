WORKSPACE ?= $(shell pwd)
include ${WORKSPACE}/.env
export $(shell sed 's/=.*//' .env)
export WORKSPACE