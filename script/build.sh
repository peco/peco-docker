#!/bin/bash

set -e

if [ -z "$PECO_DIR" ]; then
    PECO_DIR=`pwd`/../peco
fi

if [ -z "$PECO_DOCKER_IMAGE" ]; then
    PECO_DOCKER_IMAGE=peco-docker:go1.4
fi

if [ -z "$ARTIFACTS_DIR" ]; then
    docker run --rm \
        -v $PECO_DIR:/work/src/github.com/peco/peco \
        $PECO_DOCKER_IMAGE \
        /build-docker.sh
else
    docker run --rm \
        -v $ARTIFACTS_DIR:/work/artifacts \
        -v $PECO_DIR:/work/src/github.com/peco/peco \
        $PECO_DOCKER_IMAGE \
        /build-docker.sh
fi