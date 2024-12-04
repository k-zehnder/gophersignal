#!/bin/sh

# Dynamically determine the host's IP address
HOST_IP=$(ip route | grep default | awk '{print $3}')
export MYSQL_HOST=$HOST_IP
echo "Determined host IP: $MYSQL_HOST"

echo "Starting database initialization..."

# Wait for MySQL to be ready
while ! mysqladmin ping -h "$MYSQL_HOST" --silent; do
    echo "Waiting for MySQL..."
    sleep 5
done

echo "MySQL is ready."

# Start the main application
echo "Starting Go application..."

if [ "$GO_ENV" = "development" ]; then
    exec go run main.go
else
    exec ./main
fi

