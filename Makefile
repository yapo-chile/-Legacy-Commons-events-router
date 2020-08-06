include scripts/commands/vars.mk

export BRANCH ?= $(shell git branch | sed -n 's/^\* //p')
export COMMIT_DATE_UTC ?= $(shell TZ=UTC git show --quiet --date='format-local:%Y%m%d_%H%M%S' --format="%cd")

export DOCKER_TAG ?= $(shell echo ${BRANCH} | tr '[:upper:]' '[:lower:]' | sed 's,/,_,g')
export CHART_DIR ?= k8s/${APPNAME}

## Run tests and generate quality reports
test:
	@scripts/commands/test.sh

## Run tests and output coverage reports
cover:
	@scripts/commands/test_cover.sh cli

## Run tests and open report on default web browser
coverhtml:
	@scripts/commands/test_cover.sh html

## Run gometalinter and output report as text
checkstyle:
	@scripts/commands/test_style.sh display

## Install golang system level dependencies
setup:
	@scripts/commands/setup.sh
	
## Compile and build the executable file for pact tests
pact-build:
	scripts/commands/pact-build.sh

## Execute pact tests
pact-test: pact-build
	scripts/commands/pact-test.sh
	
## Compile the code
build:
	@scripts/commands/build.sh

## Upload helm charts for deploying on k8s
helm-publish:
	@echo "Publishing helm package to Artifactory"
	helm lint ${CHART_DIR}
	helm package ${CHART_DIR}
	jfrog rt u "*.tgz" "helm-local/yapo/" || true

## Execute the service
run:
	@env APP_PORT=${SERVICE_PORT} ./${APPNAME}

## Compile and start the service
start: build run

## Compile and start the service using docker
docker-start: build docker-build docker-compose-up info

## Stop docker containers
docker-stop: docker-compose-down

## Setup a new service repository based on events-router
clone:
	@scripts/commands/clone.sh

## Run gofmt to reindent source
fix-format:
	@scripts/commands/fix-format.sh

## Display basic service info
info:
	@echo "Service: ${APPNAME}"
	@echo "Images from latest commit:"
	@echo "- ${DOCKER_IMAGE}:${DOCKER_TAG}"
	@echo "- ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}"
	@echo "API Base URL: ${BASE_URL}"
	@echo "Healthcheck: curl ${BASE_URL}/api/v1/healthcheck"

include docs.mk
include docker.mk
include help.mk
