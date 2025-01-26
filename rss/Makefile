.PHONY: build
build:
	@echo "Building RSS Docker image..."
	docker build -t $(RSS_IMAGE_TAG) .

.PHONY: push
push:
	@echo "Pushing RSS Docker image..."
	docker push $(RSS_IMAGE_TAG)

.PHONY: run
run:
	@echo "Starting RSS service..."
	docker compose up -d rss

.PHONY: test
test:
	@echo "Running Rust tests with cargo..."
	cargo test

.PHONY: pull
pull:
	@echo "Pulling RSS Docker image..."
	docker pull $(RSS_IMAGE_TAG)
