#! /bin/bash

set -e

docker run -it \
--name golang-gorm-chi-postgres \
-d \
-p 6876:5432 \
-e POSTGRES_USER=user \
-e POSTGRES_PASSWORD=pass \
-e POSTGRES_DB=dev_db \
postgres:14