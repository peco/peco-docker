#!/bin/bash

set -e

if [ -z "$PECO_DIR" ]; then
    PECO_DIR=`pwd`/../peco
fi

if [ -z "$PECO_VERSION" ]; then
    echo "PECO_VERSION must be specified"
    exit 1
fi

if [ -z "$GITHUB_TOKEN_FILE" ]; then
    echo "GITHUB_TOKEN_FILE must be specified"
    exit 1
fi

if [ -z "$PECO_DOCKER_IMAGE"]; then
    PECO_DOCKER_IMAGE=peco-docker:go1.4
fi

docker run --rm \
    -v $PECO_DIR:/work/src/github.com/peco/peco \
    -e PECO_VERSION=$PECO_VERSION \
    -e GITHUB_USERNAME=peco \
    -e GITHUB_TOKEN=`cat $GITHUB_TOKEN_FILE` \
    $PECO_DOCKER_IMAGE \
    /release-docker.sh