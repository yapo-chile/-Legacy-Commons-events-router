#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

########### CODE ##############

#Build code again now for docker platform
echoHeader "Building code for docker platform"

set +e

echoTitle "Building docker image for ${DOCKER_IMAGE}"
echo "GIT BRANCH: ${GIT_BRANCH}"
echo "GIT COMMIT: ${GIT_COMMIT}"
echo "GIT COMMIT DATE: ${GIT_COMMIT_DATE}"
echo "BUILD CREATOR: ${BUILD_CREATOR}"
echo "IMAGE NAME: ${DOCKER_IMAGE}:${DOCKER_TAG}"
echo "IMAGE NAME: ${DOCKER_IMAGE}:${COMMIT_DATE_UTC}"

DOCKER_ARGS=" \
    -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
    -t ${DOCKER_IMAGE}:${COMMIT_DATE_UTC} \
    --build-arg GIT_BRANCH="$GIT_BRANCH" \
    --build-arg GIT_COMMIT="$GIT_COMMIT" \
    --build-arg GIT_COMMIT_DATE="$GIT_COMMIT_DATE" \
    --build-arg BUILD_CREATOR="$BUILD_CREATOR" \
    --build-arg APPNAME="$APPNAME" \
    -f docker/dockerfile \
    ."

echo "args: ${DOCKER_ARGS}"
set -x
docker build ${DOCKER_ARGS}
set +x

echoTitle "Build done"
