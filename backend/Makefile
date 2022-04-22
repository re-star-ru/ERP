include ../.dev.env
export

dev:
	go run ./cmd/proxy

meili:
	docker run -d -it --rm \
    -p 7700:7700 \
    -v $(shell pwd)/data.ms:/data.ms \
    getmeili/meilisearch:latest

.PHONY: backend
