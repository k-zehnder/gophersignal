#!/bin/sh

echo "Starting database initialization..."

# Wait for MySQL to be ready
while ! mysqladmin ping -h "$MYSQL_HOST" --silent; do
    echo "Waiting for MySQL..."
    sleep 5
done

echo "MySQL is ready."

# Check if the database exists and create it if it does not
echo "Creating database if it doesn't exist..."
mysql -h $MYSQL_HOST -u root -p"$MYSQL_ROOT_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE"

# Apply the schema.sql (idempotent operation)
echo "Applying schema..."
mysql -h $MYSQL_HOST -u root -p"$MYSQL_ROOT_PASSWORD" $MYSQL_DATABASE < /app/schema.sql

echo "Database initialization completed."

# Start the main application
echo "Starting Go application..."
exec go run main.go
