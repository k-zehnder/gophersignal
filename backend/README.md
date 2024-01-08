# backend

### Backend (Reminders & Edge Cases):
```markdown
1. **Docker Setup**: Run `docker-compose up -d --build --remove-orphans` after a reboot to set up Docker.

2. **MYSQL_DSN Host**: To find the MySQL host for `MYSQL_DSN` in production, use: `docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' CONTAINER_ID`.
