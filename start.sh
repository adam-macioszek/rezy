#!/bin/sh

set -e

echo "running database migrations"

/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

exec "$@"