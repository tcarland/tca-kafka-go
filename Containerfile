ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION} as build
LABEL author="Timothy C. Arland <tcarland at gmail dot com>"

WORKDIR /tcarland-kafka-go

COPY . .

RUN cd kafka && go build 
RUN go test ./utils/

ENTRYPOINT ["/usr/bin/tini", "--"]