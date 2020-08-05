#!/usr/bin/env bash

# Include colors.sh
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/colors.sh"

echoHeader "Running dependencies script"

set -e
# List of tools used for testing, validation, and report generation
tools=(
    github.com/jstemmer/go-junit-report
    github.com/axw/gocov/gocov
    github.com/AlekSi/gocov-xml
    github.com/Masterminds/glide
)

 # install librdkafka
if [ "$(id -nu)" != "root" ]; then
    echoTitle "Installing librdkafka"
    git clone https://github.com/edenhill/librdkafka.git
    cd librdkafka
    ./configure --prefix /usr
    make
    sudo make install
    cd ..
    rm -rf librdkafka
fi
# install librdkafka
if [ "$(id -nu)" != "root" ]; then
    echoTitle "Installing librdkafka"
    rm -rf librdkafka
    git clone https://github.com/edenhill/librdkafka.git
    cd librdkafka
    ./configure --prefix /usr/local
    make
    sudo make install
    cd ..
    rm -rf librdkafka
fi

echoTitle "Installing the sneaky golangci-lint"
GO111MODULE=on mod init && go get -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0

echoTitle "Installing missing tools"
# Install missed tools
for tool in ${tools[@]}; do
	GO111MODULE=off go get -u -v ${tool}
done

echoTitle "Installing Glide dependencies"
glide install

set +e
