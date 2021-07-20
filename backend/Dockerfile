FROM golang:alpine

RUN mkdir /files
COPY . /files
WORKDIR /files
RUN go build ./cmd/apiserver


ENTRYPOINT ["/files/apiserver"]