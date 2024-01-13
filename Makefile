VERSION ?= latest
FRONTEND_TAG = 'frontend'
BACKEND_TAG = 'backend'
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = 'kjzehnder3/gophersignal'

.PHONY: all build-frontend build-backend run stop down test docker_push_frontend docker_push_backend

# Default target for building and running everything
all: build-frontend build-backend run

# Build frontend Docker image
build-frontend:
	cd frontend && docker build -t $(DOCKERHUB_REPO):$(FRONTEND_TAG) .

# Build backend Docker image
build-backend:
	cd backend && docker build -t $(DOCKERHUB_REPO):$(BACKEND_TAG) .

# Start all services using Docker Compose
run:
	$(DOCKER_COMPOSE) up -d

# Stop all services without removing them
stop:
	$(DOCKER_COMPOSE) stop

# Stop and remove all running services
down:
	$(DOCKER_COMPOSE) down

# Run tests (TODO: frontend tests)
test:
	# Backend tests
	cd backend && go test -v -cover ./...

# Push frontend Docker image to Docker Hub
docker_push_frontend:
	docker push $(DOCKERHUB_REPO):$(FRONTEND_TAG)

# Push backend Docker image to Docker Hub
docker_push_backend:
	docker push $(DOCKERHUB_REPO):$(BACKEND_TAG)
