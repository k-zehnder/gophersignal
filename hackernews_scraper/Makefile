.PHONY: build
build:
	@echo "Building HackerNews scraper Docker image..."
	docker build -t $(HACKERNEWS_SCRAPER_IMAGE_TAG) .

.PHONY: push
push:
	@echo "Pushing HackerNews scraper Docker image..."
	docker push $(HACKERNEWS_SCRAPER_IMAGE_TAG)

.PHONY: run
run:
	@echo "Starting HackerNews scraper service..."
	docker compose up -d hackernews_scraper

.PHONY: scrape
run:
	@echo "Starting HackerNews scraper service..."
	npm run start	