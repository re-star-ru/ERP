FROM ghcr.io/umputun/baseimage/buildgo:latest as build

WORKDIR /build
ADD . /build

RUN ls /build
RUN go build -o app ./cmd/proxy

#################################################################

FROM ghcr.io/umputun/baseimage/app:lastest

COPY --from=build /build/app /srv/app
WORKDIR /srv

CMD ["/srv/app"]