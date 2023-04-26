#! /bin/bash

set -e

swag init

go build -o bin/ ./...
