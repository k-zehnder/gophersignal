.PHONY: build run test

build:
	docker build -t $(BACKEND_TAG) .

run:
	docker-compose up -d backend

test:
	go test -v -cover ./...
