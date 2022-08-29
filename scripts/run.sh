#!/bin/sh

set -e

go mod tidy
go run main.go "$@"