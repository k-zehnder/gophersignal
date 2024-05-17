#!/bin/sh

# Conditionally initialize the database if the INIT_DB variable is true
if [ "$INIT_DB" = "true" ]; then
  echo "Initializing database..."
  # Use the MySQL client to execute the schema.sql file
  mysql -h "$MYSQL_HOST" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < /app/schema.sql
  echo "Database initialization completed."
fi

# Start the main Go application
exec go run main.go
