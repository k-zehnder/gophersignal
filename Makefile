VERSION ?= latest
DOCKER_COMPOSE = docker-compose
DOCKERHUB_REPO = kjzehnder3/gophersignal
FRONTEND_TAG = $(DOCKERHUB_REPO):frontend-$(VERSION)
BACKEND_TAG = $(DOCKERHUB_REPO):backend-$(VERSION)

export FRONTEND_TAG BACKEND_TAG DOCKERHUB_REPO

# Common Targets
.PHONY: all
all: build_frontend build_backend test_backend

.PHONY: all_push
all_push: build_frontend build_backend test_backend docker_push_backend docker_push_frontend

.PHONY: run
run:
	docker-compose up -d
	@echo "Done!"

.PHONY: deploy_pull
deploy_pull: 
	@echo "Pulling frontend and backend, then deploying..."
	docker pull frontend
	docker pull backend
	docker-compose up -d
	@echo "Frontend and Backend running"

.PHONY: dev
dev:
	@echo "Backend services are starting. Please wait..."
	$(DOCKER_COMPOSE) -f docker-compose-dev.yml up -d --build
	@echo "Frontend server is initializing and may take a minute to become available at localhost:3000."
	@echo "Development environment setup complete. If the backend doesn't start successfully, you might need to run 'make dev' again."

.PHONY: dev_data
dev_data:
	@echo "Setting up development data..."
	@echo "Example command to set up development data with custom values:"
	@echo "make dev_data HUGGING_FACE_API_KEY=your_api_key_here SCRAPER_MYSQL_DSN=mysql_dsn_here"
	@HUGGING_FACE_API_KEY=$(HUGGING_FACE_API_KEY) \
	SCRAPER_MYSQL_DSN="$(SCRAPER_MYSQL_DSN)" \
	cd backend && make scrape
	@echo "Data scraping complete."
	@echo "Hugging Face API key has been set."
	@HUGGING_FACE_API_KEY=$(HUGGING_FACE_API_KEY) \
	cd backend && make hfsummarize
	@echo "Summarization complete."

.PHONY: scrape
scrape:
	@echo "Scraping data..."
	cd backend && make scrape
	@echo "Scraping complete"

.PHONY: openaisummarize
openaisummarize:
	@echo "Summarizing with OpenAI..."
	cd backend && make openaisummarize
	@echo "OpenAI summarization complete"

.PHONY: hfsummarize
hfsummarize:
	@echo "Summarizing with HF..."
	cd backend && make hfsummarize
	@echo "HF summarization complete"

# Frontend Section
.PHONY: build_frontend
build_frontend:
	@echo "Building frontend..."
	cd frontend && make build_frontend
	@echo "Frontend built successfully"

.PHONY: run_frontend
run_frontend:
	@echo "Starting frontend..."
	docker run -d --name frontend $(FRONTEND_TAG)
	@echo "Frontend running"

.PHONY: docker_push_frontend
docker_push_frontend:
	@echo "Pushing frontend image..."
	docker push $(FRONTEND_TAG)
	@echo "Frontend image pushed successfully"

.PHONY: pull_frontend
pull_frontend:
	@echo "Pulling frontend image..."
	docker pull $(FRONTEND_TAG)
	@echo "Frontend image pulled successfully"

.PHONY: deploy_frontend
deploy_frontend:
	@echo "Deploying frontend..."
	docker-compose up -d frontend
	@echo "Frontend running"

# Backend Section
.PHONY: build_backend
build_backend:
	@echo "Building backend..."
	cd backend && make build_backend
	@echo "Backend built successfully"

.PHONY: run_backend
run_backend:
	@echo "Starting backend..."
	cd backend && make run_backend
	@echo "Backend running"

.PHONY: docker_push_backend
docker_push_backend:
	@echo "Pushing backend image..."
	cd backend && docker push $(BACKEND_TAG)
	@echo "Backend image pushed successfully"

.PHONY: pull_backend
pull_backend:
	@echo "Pulling backend image..."
	cd backend && docker pull $(BACKEND_TAG)
	@echo "Backend image pulled successfully"

.PHONY: deploy_backend
deploy_backend:
	@echo "Deploying backend..."
	docker-compose up -d backend
	@echo "Backend running"

.PHONY: test_backend
test_backend:
	@echo "Running backend tests..."
	cd backend && make test_backend
	@echo "Backend tests completed"
