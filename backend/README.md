# backend

**Note**: To find the MySQL host for `MYSQL_DSN` in production, use: `docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' CONTAINER_ID`.