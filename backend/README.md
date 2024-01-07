# backend

### Backend (Reminders & Edge Cases):
```markdown
1. **Docker Setup**: Run `docker-compose up -d --build --remove-orphans` after a reboot to set up Docker.

2. **MYSQL_DSN Host**: To find the MySQL host for `MYSQL_DSN` in production, use: `docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' CONTAINER_ID`.

3. **MYSQL_DSN Issues**: Post-reboot, if there are connectivity issues, switch `MYSQL_DSN` between `XXX.XXX.XXX:3333` and `mysql:3333`.

4. **.env File Duplication**: Resolve duplication of `.env` files in both `/` and `/backend`.

5. **Go Test Coverage**: Run `go test -coverprofile=coverage.out ./...` for coverage, then `go tool cover -html=coverage.out -o coverage.html` to generate an HTML report. View it with `open coverage.html`.

6. **Pull Docker Images Manually**: If needed, pull the necessary Docker images manually using:
   - `docker pull node:latest`
   - `docker pull nginx:latest`
   - `docker pull alpine:latest`
```
