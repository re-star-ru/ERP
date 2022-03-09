include .env
export

frontend:
	cd frontend && npm i

backend:
	cd backend && go run ./cmd/proxy

backend-dev:
	cd backend && HOST=8100 go run ./cmd/proxy

.PHONY: frontend backend