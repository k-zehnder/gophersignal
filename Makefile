VERSION ?= latest
DOCKER_COMPOSE_CMD = docker-compose
DOCKERHUB_REPO = kjzehnder3/gophersignal
FRONTEND_IMAGE_TAG = $(DOCKERHUB_REPO):frontend-$(VERSION)
BACKEND_IMAGE_TAG = $(DOCKERHUB_REPO):backend-$(VERSION)

export FRONTEND_IMAGE_TAG BACKEND_IMAGE_TAG DOCKERHUB_REPO

# Common Targets
.PHONY: all
all: build_frontend build_backend test_backend
	@echo "All components built and tested!"

.PHONY: push_all
push_all: build_frontend build_backend test_backend push_backend push_frontend

.PHONY: run_services
run_services:
	$(DOCKER_COMPOSE_CMD) up -d
	@echo "Frontend: http://localhost:3000"
	@echo "Backend: http://localhost:8080"
	@echo "Swagger Docs: http://localhost:8080/swagger/index.html"

.PHONY: deploy
deploy: 
	@echo "Deploying frontend and backend..."
	docker pull $(FRONTEND_IMAGE_TAG)
	docker pull $(BACKEND_IMAGE_TAG)
	$(DOCKER_COMPOSE_CMD) up -d
	@echo "Services deployed successfully."

# Development Environment
.PHONY: dev
dev:
	@echo "Starting backend services..."
	$(DOCKER_COMPOSE_CMD) -f docker-compose.dev.yml up -d --build
	@echo "Frontend initializing at http://localhost:3000"
	@echo "Development environment ready."

.PHONY: setup_dev
setup_dev:
	@echo "Setting up development data..."
	@echo "Use 'make setup_dev_data HUGGING_FACE_API_KEY=<key> MYSQL_DSN=<dsn>' to customize."
	HUGGING_FACE_API_KEY=$(HUGGING_FACE_API_KEY) MYSQL_DSN="$(MYSQL_DSN)" cd backend && make scrape
	@echo "Data scraping and setup complete."

.PHONY: scrape_data
scrape_data:
	@echo "Scraping data..."
	cd backend && make scrape_data
	@echo "Data scraping complete."

# Summarization
.PHONY: summarize_openai
summarize_openai:
	@echo "Summarizing with OpenAI (ChatGPT)..."
	cd backend && make summarize_openai
	@echo "OpenAI summarization complete."

.PHONY: summarize_huggingface
summarize_huggingface:
	@echo "Summarizing with Hugging Face..."
	cd backend && make summarize_huggingface
	@echo "Hugging Face summarization complete."

# Frontend Management
.PHONY: build_frontend
build_frontend:
	@echo "Building frontend..."
	cd frontend && make build_frontend
	@echo "Frontend build complete."

.PHONY: start_frontend
start_frontend:
	@echo "Starting frontend..."
	docker run -d --name frontend $(FRONTEND_IMAGE_TAG)
	@echo "Frontend service started."

.PHONY: push_frontend
push_frontend:
	@echo "Pushing frontend image..."
	docker push $(FRONTEND_IMAGE_TAG)
	@echo "Frontend image pushed."

.PHONY: pull_frontend
pull_frontend:
	@echo "Pulling frontend image..."
	docker pull $(FRONTEND_IMAGE_TAG)
	@echo "Frontend image pulled."

.PHONY: deploy_frontend
deploy_frontend:
	@echo "Deploying frontend..."
	$(DOCKER_COMPOSE_CMD) up -d frontend
	@echo "Frontend deployed."

# Backend Management
.PHONY: build_backend
build_backend:
	@echo "Building backend..."
	cd backend && make build_backend
	@echo "Backend build complete."

.PHONY: start_backend
start_backend:
	@echo "Starting backend..."
	cd backend && make start
	@echo "Backend service started."

.PHONY: push_backend
push_backend:
	@echo "Pushing backend image..."
	cd backend && docker push $(BACKEND_IMAGE_TAG)
	@echo "Backend image pushed."

.PHONY: pull_backend
pull_backend:
	@echo "Pulling backend image..."
	cd backend && docker pull $(BACKEND_IMAGE_TAG)
	@echo "Backend image pulled."

.PHONY: deploy_backend
deploy_backend:
	@echo "Deploying backend..."
	$(DOCKER_COMPOSE_CMD) up -d backend
	@echo "Backend deployed."

.PHONY: test_backend
test_backend:
	@echo "Running backend tests..."
	cd backend && make test_backend
	@echo "Backend tests completed."

.PHONY: test_coverage
test_coverage:
	@echo "Running Go tests with coverage in backend..."
	cd backend && go test -coverprofile=coverage.out ./...
	@echo "Generating coverage report in backend..."
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	cd backend && open coverage.html
