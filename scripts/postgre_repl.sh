#! /bin/bash
set -e
export PGPASSWORD=pass
psql -h localhost -d dev_db -U user -p 6876