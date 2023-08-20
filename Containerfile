ARG GO_VERSION=1.19

FROM golang:${GO_VERSION} as build
LABEL author="Timothy C. Arland <tcarland at gmail dot com>"

WORKDIR /tcarland-kafka-go

COPY . .

RUN cd kafka && go build 
RUN go test ./utils/ -v
RUN cd utils && go build

ENTRYPOINT ["/usr/bin/tini", "--"]