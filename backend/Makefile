include ../.dev.env
export

# rdocker: 
#     docker run --rm \
#         -p 8100:8100 \
#         -e MINIO_ACCESS_KEY=ryR3bGLYsGXV \
#         -e MINIO_SECRET_KEY=AoRqxFTVWs3U \
#         -e ADDR=0.0.0.0:8100 \

dev:
	go run ./cmd/proxy

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

.PHONY: backend rdocker dev
