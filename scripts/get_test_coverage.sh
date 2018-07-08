#!/usr/bin/env bash
set -e
go test ./sla-calculator -covermode=count -coverprofile=cover.out
go tool cover -func=cover.out
$HOME/gopath/bin/goveralls -coverprofile=cover.out -service=travis-ci