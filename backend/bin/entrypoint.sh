#!/bin/sh
set -euo pipefail

# Dynamically determine the host's IP address if MYSQL_HOST not set
if [ -z "${MYSQL_HOST:-}" ]; then
  HOST_IP="$(ip -4 route show default | awk '/default/ {print $3}')"
  export MYSQL_HOST="$HOST_IP"
fi
: "${MYSQL_PORT:=3306}"
: "${MYSQL_DATABASE:?gophersignal}"
: "${MYSQL_USER:?user}"
: "${MYSQL_PASSWORD:?}"
: "${GO_ENV:=production}"

echo "Determined host IP: $MYSQL_HOST"
echo "Starting database initialization..."

# Wait for MySQL to be ready
export MYSQL_PWD="${MYSQL_PASSWORD}"
until mysqladmin --protocol=tcp -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
  -u "$MYSQL_USER" ping >/dev/null 2>&1; do
  echo "Waiting for MySQL..."
  sleep 5
done
echo "MySQL is ready."

# Apply the schema.sql (idempotent) as APP USER (no root)
if [ -f /app/schema.sql ]; then
  echo "Applying schema..."
  mysql --protocol=tcp -h "$MYSQL_HOST" -P "$MYSQL_PORT" \
    -u "$MYSQL_USER" "$MYSQL_DATABASE" < /app/schema.sql
  echo "Schema applied."
else
  echo "No /app/schema.sql found; skipping schema apply."
fi

unset MYSQL_PWD
echo "Database initialization completed."

# Start the main application
echo "Starting Go application..."
if [ "${GO_ENV}" = "development" ]; then
  exec go run main.go
else
  exec ./main
fi
