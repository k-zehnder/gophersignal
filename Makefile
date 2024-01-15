VERSION ?= latest
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = kjzehnder3/gophersignal
FRONTEND_TAG = $(DOCKERHUB_REPO):frontend-$(VERSION)
BACKEND_TAG = $(DOCKERHUB_REPO):backend-$(VERSION)

export FRONTEND_TAG BACKEND_TAG DOCKERHUB_REPO

.PHONY: all build_frontend build_backend run stop down test docker_push_frontend docker_push_backend all_push

all: build_frontend build_backend run

build_frontend:
	cd frontend && make build

build_backend:
	cd backend && make build

run:
	$(DOCKER_COMPOSE) up -d

stop:
	$(DOCKER_COMPOSE) stop

down:
	$(DOCKER_COMPOSE) down

test:
	cd backend && make test

docker_push_frontend:
	docker push $(FRONTEND_TAG)

docker_push_backend:
	docker push $(BACKEND_TAG)

all_push: build_frontend build_backend test docker_push_frontend docker_push_backend
