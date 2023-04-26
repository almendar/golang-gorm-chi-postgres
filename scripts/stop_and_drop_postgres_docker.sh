#! /bin/bash

set -e

docker stop golang-gorm-postgres
docker rm golang-gorm-postgres