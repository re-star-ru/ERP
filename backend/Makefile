include ../.dev.env
export

export COMPOSE_PROJECT_NAME=oprox
compose:
	docker compose up

gen:
	go generate ./...

dev:
	CGO_ENABLED=0 go run ./cmd/proxy

dent:
	CGO_ENABLED=0 go run ./ent/start

meili:
	docker run -d -it --rm \
	-p 7700:7700 \
	-v $(shell pwd)/data.ms:/data.ms \
	-v $(shell pwd)/dumps:/dumps \
	getmeili/meilisearch:latest

indexUpd:
	xh post localhost:7700/indexes/items/settings/typo-tolerance --raw '{ "disableOnAttributes": [], "enabled": false }'

clear:
	xh delete localhost:7700/indexes/items/documents

filter-purple:
	 xh POST -j localhost:7700/indexes/items/search filter="lastExpance < 1557131207"

.PHONY: backend rdocker dev compose
