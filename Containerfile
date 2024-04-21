ARG GO_VERSION=1.22

FROM golang:${GO_VERSION} as build
LABEL author="Timothy C. Arland <tcarland at gmail dot com>"

WORKDIR /tca-kafka-go

COPY . .

RUN cd kafka && go build 
RUN cd utils && go build
RUN go test ./utils/ -v

ENTRYPOINT ["/usr/bin/tini", "--"]