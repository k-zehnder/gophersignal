VERSION ?= latest
FRONTEND_TAG = $(DOCKERHUB_REPO):frontend-$(VERSION)
BACKEND_TAG = $(DOCKERHUB_REPO):backend-$(VERSION)
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = kjzehnder3/gophersignal

.PHONY: all build-frontend build-backend run stop down test docker_push_frontend docker_push_backend

all: build-frontend build-backend run

build-frontend:
	cd frontend && docker build -t $(FRONTEND_TAG) .

build-backend:
	cd backend && docker build -t $(BACKEND_TAG) .

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
	docker push $(FRONTEND_TAG)

docker_push_backend:
	docker push $(BACKEND_TAG)
