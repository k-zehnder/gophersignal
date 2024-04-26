.PHONY: build_frontend
build_frontend:
	docker build -t $(FRONTEND_IMAGE_TAG) .

.PHONY: run_frontend
run_frontend:
	docker-compose up -d frontend

.PHONY: dev
dev:
	npm run dev

.PHONY: test_frontend
test_frontend:
	npm i jest
	npm i jest-environment-jsdom
	npm run test
