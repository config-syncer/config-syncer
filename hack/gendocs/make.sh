#!/usr/bin/env bash

pushd $GOPATH/src/github.com/appscode/kubed/hack/gendocs
go run main.go
popd
