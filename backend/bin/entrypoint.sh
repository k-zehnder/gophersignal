#!/bin/sh
set -euo pipefail

# Prefer explicit env; default to host.docker.internal on Linux
: "${MYSQL_HOST:=host.docker.internal}"
: "${MYSQL_PORT:=3306}"
: "${MYSQL_DATABASE:=gophersignal}"
: "${MYSQL_USER:=user}"
: "${MYSQL_PASSWORD:=password}"
: "${GO_ENV:=production}"

echo "Using MySQL ${MYSQL_HOST}:${MYSQL_PORT} db=${MYSQL_DATABASE} user=${MYSQL_USER}"
echo "Starting database initialization..."

export MYSQL_PWD="${MYSQL_PASSWORD}"
until mysqladmin --protocol=tcp -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" ping >/dev/null 2>&1; do
  echo "Waiting for MySQL..."
  sleep 5
done
echo "MySQL is ready."

if [ -f /app/schema.sql ]; then
  echo "Applying schema..."
  mysql --protocol=tcp -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" "$MYSQL_DATABASE" < /app/schema.sql
  echo "Schema applied."
else
  echo "No /app/schema.sql found; skipping schema apply."
fi

unset MYSQL_PWD
echo "Database initialization completed."
echo "Starting Go application..."
if [ "${GO_ENV}" = "development" ]; then
  exec go run main.go
else
  exec ./main
fi
