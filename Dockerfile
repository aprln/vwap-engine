FROM golang:1.19-alpine
RUN apk add build-base
ADD .. /vwap-engine
WORKDIR /vwap-engine

RUN go mod tidy
RUN go mod vendor
RUN go build -race -o bin/vwap
