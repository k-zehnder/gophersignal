.PHONY: build_frontend
build_frontend:
	docker build -t $(FRONTEND_IMAGE_TAG) .

.PHONY: run_frontend
run_frontend:
	docker-compose up -d frontend

.PHONY: dev
dev:
	npm install && npm run dev
