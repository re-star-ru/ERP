frontend:
	cd frontend && npm i

backend:
	cd backend && go run ./cmd/proxy

.PHONY: frontend backend