VERSION ?= latest
FRONTEND_TAG = $(DOCKERHUB_REPO):frontend-$(VERSION)
BACKEND_TAG = $(DOCKERHUB_REPO):backend-$(VERSION)
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = kjzehnder3/gophersignal

export FRONTEND_TAG BACKEND_TAG DOCKERHUB_REPO

.PHONY: all build-frontend build-backend run stop down test docker_push_frontend docker_push_backend all-push

all: build-frontend build-backend run

build-frontend:
	cd frontend && make build

build-backend:
	cd backend && make build

run:
	$(DOCKER_COMPOSE) up -d

stop:
	$(DOCKER_COMPOSE) stop

down:
	$(DOCKER_COMPOSE) down

test:
	cd backend && go test -v -cover ./...

docker_push_frontend:
	docker push $(FRONTEND_TAG)

docker_push_backend:
	docker push $(BACKEND_TAG)

all-push: build-frontend build-backend test docker_push_frontend docker_push_backend
