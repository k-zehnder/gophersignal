.PHONY: build
build:
	@echo "Building HackerNews scraper Docker image..."
	docker build -t $(HACKERNEWS_SCRAPER_IMAGE_TAG) .

.PHONY: test
test:
	@echo "Running tests for HackerNews scraper..."
	docker run --rm $(HACKERNEWS_SCRAPER_IMAGE_TAG) npm test

.PHONY: push
push:
	@echo "Pushing HackerNews scraper Docker image..."
	docker push $(HACKERNEWS_SCRAPER_IMAGE_TAG)

.PHONY: run
run:
	@echo "Starting HackerNews scraper service..."
	docker compose up -d hackernews_scraper

.PHONY: pull
pull:
	@echo "Pulling hackernews_scraper Docker image..."
	docker pull $(HACKERNEWS_SCRAPER_IMAGE_TAG)

.PHONY: scrape 
scrape:
	@echo "Running HackerNews scraper..."
	docker exec -it $(shell docker ps -qf "name=hackernews_scraper") npm run start
