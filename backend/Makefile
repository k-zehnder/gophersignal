.PHONY: build
build_backend:
	docker build -t $(BACKEND_TAG) .

.PHONY: run
run_backend:
	docker-compose up -d backend

.PHONY: test_backend
test_backend:
	go test -v -cover ./...
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: scrape
scrape:
	go run cmd/articlesaver/main.go

.PHONY: openaisummarize
openaisummarize:
	go run cmd/openaisummarizer/main.go

.PHONY: hfsummarize
hfsummarize:
	go run cmd/hfsummarizer/main.go