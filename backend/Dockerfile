FROM ghcr.io/umputun/baseimage/buildgo:latest as build

WORKDIR /build
ADD . /build

RUN ls /build
RUN go build -o app ./cmd/proxy

#################################################################

FROM ghcr.io/umputun/baseimage/app:lastest

ENV TIME_ZONE=Europe/Moscow
ENV ADDR="0.0.0.0:8080"
ENV MINIO_ENDPOINT=https://s3.re-star.ru
ENV ONEC_HOST=http://10.51.0.10:19201/sm1/hs/

ARG MINIO_ACCESS_KEY
ARG MINIO_SECRET_KEY
ARG ONEC_TOKEN

COPY --from=build /build/app /srv/app
WORKDIR /srv

CMD ["/srv/app"]