FROM node:current-alpine as backend

ARG GIT_BRANCH

ADD . /build
WORKDIR /build
#RUN apk add --no-cache --update git tzdata ca-certificates
# TODO: -D for draft build

RUN npm i
RUN npm run build

#RUN \
#    version=${GIT_BRANCH}-$(date +%Y%m%dT%H:%M:%S) && \
#    echo "version=$version" && \
#    cd cmd/app && \
#    go build -o /build/rp -ldfl ags "-X main.revision=${version} -s -w"

FROM umputun/reproxy:latest
#todo: enable reproxy SPA mod
COPY --from=backend /build/dist/spa /public
ENV ASSETS_LOCATION=/public
ENV ASSETS_SPA=true
WORKDIR /public
