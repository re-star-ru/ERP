include ../.dev.env
export

dev:
	go run ./cmd/proxy

MINIO_ACCESS_KEY=ryR3bGLYsGXV
MINIO_SECRET_KEY=AoRqxFTVWs3U

rdocker: 
	docker run --rm \
		-p 8100:8100 \
		-e MINIO_ACCESS_KEY=ryR3bGLYsGXV \
		-e MINIO_SECRET_KEY=AoRqxFTVWs3U \
		-e ADDR=0.0.0.0:8100 \
		a564da9db929 

meili:
	docker run -d -it --rm \
    -p 7700:7700 \
    -v $(shell pwd)/data.ms:/data.ms \
    getmeili/meilisearch:latest

.PHONY: backend rdocker
