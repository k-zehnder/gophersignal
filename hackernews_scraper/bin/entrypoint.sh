#!/bin/sh

# Wait for the backend service to be ready
until nc -z backend 8080; do
  echo "Waiting for backend service to be ready..."
  sleep 10
done

echo "Backend service is ready. Starting hackernews scraper..."

# Run the scraper once if INIT_DB is true
if [ "$INIT_DB" = "true" ]; then
  npm run start
fi

# Keep the container running
tail -f /dev/null
