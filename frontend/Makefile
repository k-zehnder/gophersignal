.PHONY: build
build_frontend:
	docker build -t $(FRONTEND_TAG) .

.PHONY: run
run_frontend:
	docker-compose up -d frontend

.PHONY: dev
dev:
	npm install && npm run dev
