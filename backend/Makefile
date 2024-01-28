.PHONY: build_backend
build_backend:
	docker build -t $(BACKEND_TAG) .

.PHONY: run_backend
run_backend:
	docker-compose up -d backend

.PHONY: test_backend
test_backend:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: scrape_data
scrape_data:
	go run cmd/articlesaver/main.go

.PHONY: summarize_openai
summarize_openai:
	go run cmd/openaisummarizer/main.go

.PHONY: summarize_huggingface
summarize_huggingface:
	go run cmd/huggingfacesummarizer/main.go

.PHONY: test_coverage
test_coverage:
	@echo "Running Go tests with coverage..."
	go test -coverprofile=coverage.out ./...
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Opening coverage report..."
	open coverage.html
