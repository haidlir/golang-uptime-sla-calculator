#!/usr/bin/env bash
set -e
go get ./...
go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls