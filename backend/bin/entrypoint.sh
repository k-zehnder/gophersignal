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
mysql -h $MYSQL_HOST -u root -p"$MYSQL_ROOT_PASSWORD" -e "
    CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE;
    USE $MYSQL_DATABASE;
    CREATE TABLE IF NOT EXISTS articles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        link VARCHAR(512) NOT NULL,
        content TEXT,
        summary VARCHAR(2000),
        source VARCHAR(100) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL
    );
"

echo "Database initialization completed."

# Start the main application
echo "Starting Go application..."

if [ "$ENVIRONMENT" = "development" ]; then
    exec go run main.go
else
    exec ./main
fi
