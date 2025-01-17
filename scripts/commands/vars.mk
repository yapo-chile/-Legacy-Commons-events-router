#!/usr/bin/env bash
export GO_FILES = $(shell find . -iname '*.go' -type f | grep -v vendor | grep -v pact) # All the .go files, excluding vendor/ and pact/
GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

# GIT variables
export GIT_COMMIT=$(shell git rev-parse HEAD)
export GIT_COMMIT_DATE=$(shell TZ="America/Santiago" git show --quiet --date='format-local:%d-%m-%Y_%H:%M:%S' --format="%cd")
export BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)

# REPORT_ARTIFACTS should be in sync with `RegexpFilePathMatcher` in
# `reports-publisher/config.json`
export REPORT_ARTIFACTS=reports

# APP variables
# This variables are for the use of your microservice. This variables must be updated each time you are creating a new microservice
export APPNAME=events-router
export SERVICE_PORT=8080
export SERVICE_HOST=:localhost
export SERVER_ROOT=${PWD}
export BASE_URL="http://"${SERVICE_HOST}":"${SERVICE_PORT}"
export MAIN_FILE=cmd/${APPNAME}/main.go
export LOGGER_SYSLOG_ENABLED=false
export LOGGER_STDLOG_ENABLED=true
export LOGGER_LOG_LEVEL=0

# Pact test variables
export PACT_MAIN_FILE=cmd/${APPNAME}/main.go
export PACT_BINARY=${APPNAME}-pact
export PACT_DIRECTORY=pact
export PACT_TEST_ENABLED=false

# DOCKER variables
export DOCKER_REGISTRY=containers.mpi-internal.com
export DOCKER_IMAGE=${DOCKER_REGISTRY}/yapo/${APPNAME}
export DOCKER_PORT=$(call genport,1)

# Documentation variables
export DOCS_DIR=docs
export DOCS_HOST=localhost:$(call genport,3)
export DOCS_PATH=github.mpi-internal.com/Yapo/${APPNAME}
export DOCS_COMMIT_MESSAGE=Generate updated documentation

# Prometheus variables
export PROMETHEUS_PORT=8877
export PROMETHEUS_ENABLED=true

# Goms Client variables
export GOMS_HEALTH_PATH=http://localhost:${SERVICE_PORT}/api/v1/healthcheck

# Circuit breaker variables
export CIRCUIT_BREAKER_FAILURE_RATIO=0.5
export CIRCUIT_BREAKER_CONSECUTIVE_FAILURE=2

# User config
export PROFILE_HOST=http://10.15.1.78:7987

#Pact broker
export PACT_BROKER_HOST=http://3.229.36.112
export PACT_BROKER_PORT=80
export PROVIDER_HOST=http://localhost
export PROVIDER_PORT=8080

export KAFKA_CONSUMER_HOST=172.21.1.95
export KAFKA_CONSUMER_PORT=9092
export KAFKA_CONSUMER_TOPICS=dev01-events-queue
export KAFKA_CONSUMER_GROUP_ID=router-consumer-dev01
export KAFKA_PRODUCER_HOST=172.21.1.95
export KAFKA_PRODUCER_PORT=9092
export ETCD_HOST=http://10.15.1.78:56146
