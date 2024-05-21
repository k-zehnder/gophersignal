VERSION ?= latest
DOCKERHUB_REPO := kjzehnder3/gophersignal
FRONTEND_IMAGE_TAG := $(DOCKERHUB_REPO)-frontend:$(VERSION)
BACKEND_IMAGE_TAG := $(DOCKERHUB_REPO)-backend:$(VERSION)
HACKERNEWS_SCRAPER_IMAGE_TAG := $(DOCKERHUB_REPO)-hackernews_scraper:$(VERSION) 

export FRONTEND_IMAGE_TAG BACKEND_IMAGE_TAG HACKERNEWS_SCRAPER_IMAGE_TAG

.PHONY: all
all: build test push

.PHONY: build
build:
	@echo "Building all components..."
	$(MAKE) -C frontend build
	$(MAKE) -C backend build
	$(MAKE) -C hackernews_scraper build

.PHONY: test
test:
	@echo "Running tests for all components..."
	$(MAKE) -C frontend test
	$(MAKE) -C backend test

.PHONY: push
push:
	@echo "Pushing all images..."
	$(MAKE) -C frontend push
	$(MAKE) -C backend push
	$(MAKE) -C hackernews_scraper push

.PHONY: deploy
deploy:
	@echo "Deploying application..."
	docker compose down
	docker pull $(FRONTEND_IMAGE_TAG)
	docker pull $(BACKEND_IMAGE_TAG)
	docker pull $(HACKERNEWS_SCRAPER_IMAGE_TAG)
	docker compose up -d --build
	@echo "Application deployed successfully."

.PHONY: dev
dev:
	@echo "Starting development environment..."
	docker compose -f docker-compose-dev.yml up -d --build
