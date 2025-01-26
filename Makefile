VERSION ?= latest
DOCKERHUB_REPO := kjzehnder3/gophersignal
FRONTEND_IMAGE_TAG := $(DOCKERHUB_REPO)-frontend:$(VERSION)
BACKEND_IMAGE_TAG := $(DOCKERHUB_REPO)-backend:$(VERSION)
HACKERNEWS_SCRAPER_IMAGE_TAG := $(DOCKERHUB_REPO)-hackernews_scraper:$(VERSION)
RSS_IMAGE_TAG := $(DOCKERHUB_REPO)-rss:$(VERSION)

export FRONTEND_IMAGE_TAG BACKEND_IMAGE_TAG HACKERNEWS_SCRAPER_IMAGE_TAG RSS_IMAGE_TAG

.PHONY: all
all: build test push

.PHONY: build
build:
	@echo "Building all components..."
	$(MAKE) -C frontend build
	$(MAKE) -C backend build
	$(MAKE) -C hackernews_scraper build
	$(MAKE) -C rss build

.PHONY: test
test:
	@echo "Running tests for all components..."
	$(MAKE) -C frontend test
	$(MAKE) -C backend test
	$(MAKE) -C hackernews_scraper test
	$(MAKE) -C rss test

.PHONY: push
push:
	@echo "Pushing all images..."
	$(MAKE) -C frontend push
	$(MAKE) -C backend push
	$(MAKE) -C hackernews_scraper push
	$(MAKE) -C rss push

.PHONY: deploy
deploy:
	@echo "Deploying application..."
	docker compose down
	$(MAKE) -C frontend pull
	$(MAKE) -C backend pull
	$(MAKE) -C hackernews_scraper pull
	$(MAKE) -C rss pull
	@echo "Building frontend..."
	cd frontend && npm install && npm run build && cd ..
	@echo "Restarting Nginx..."
	docker compose up -d
	docker compose restart nginx
	@echo "Application deployed successfully."

.PHONY: dev
dev:
	@echo "Starting development environment..."
	docker compose -f docker-compose-dev.yml up -d --build

.PHONY: scrape
scrape:
	@echo "Running HackerNews Scraper inside container..."
	docker compose run --rm hackernews_scraper npm run start
