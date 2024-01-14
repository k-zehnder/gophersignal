VERSION ?= latest
FRONTEND_TAG = 'frontend'
BACKEND_TAG = 'backend'
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = 'kjzehnder3/gophersignal'

.PHONY: all build-frontend build-backend run stop down test docker_push_frontend docker_push_backend

all: build-frontend build-backend run

build-frontend:
	cd frontend && docker build -t $(DOCKERHUB_REPO):$(FRONTEND_TAG) .

build-backend:
	cd backend && docker build -t $(DOCKERHUB_REPO):$(BACKEND_TAG) .

run:
	$(DOCKER_COMPOSE) up -d

stop:
	$(DOCKER_COMPOSE) stop

down:
	$(DOCKER_COMPOSE) down

test:
	# Backend tests
	cd backend && go test -v -cover ./...

docker_push_frontend:
	docker push $(DOCKERHUB_REPO):$(FRONTEND_TAG)

docker_push_backend:
	docker push $(DOCKERHUB_REPO):$(BACKEND_TAG)
