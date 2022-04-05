FROM ghcr.io/umputun/baseimage/buildgo:latest as build

WORKDIR /build
ADD . /build

RUN ls /build
RUN go build -o app ./cmd/proxy

#################################################################

FROM ghcr.io/umputun/baseimage/app:lastest

ENV TZ=Europe/Moscow
ENV ADDR="0.0.0.0:8080"

COPY --from=build /build/app /srv/app
WORKDIR /srv

CMD ["/srv/app"]