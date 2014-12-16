#!/bin/bash

set -e

mkdir -p /work/artifacts/snapshot
cd /work/src/github.com/peco/peco
export GOPATH=/work:/work/src/github.com/peco/peco
go run build/make.go build /work/artifacts