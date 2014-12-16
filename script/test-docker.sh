#!/bin/bash

set -e

cd /work/src/github.com/peco/peco
export GOPATH=/work:/work/src/github.com/peco/peco

go run build/make.go deps
go test -v .