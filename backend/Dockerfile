FROM ghcr.io/umputun/baseimage/buildgo:latest as build

WORKDIR /build
ADD . /build

RUN ls /build
RUN go build -o app ./cmd/proxy

#################################################################

FROM umputun/baseimage:scratch-latest
ENV TZ=Europe/Moscow
ENV ADDR="0.0.0.0:8080"

ENV MINIO_ENDPOINT=s3.re-star.ru
ENV ONEC_HOST=http://10.51.0.10:19201/sm1/hs/

ARG MINIO_ACCESS_KEY
ARG MINIO_SECRET_KEY
ARG ONEC_TOKEN

COPY --from=build /build/app /srv/app
COPY --from=build /build/assets/fonts /srv/assets/fonts
COPY --from=build /build/web /srv/web

WORKDIR /srv
CMD ["/srv/app"]