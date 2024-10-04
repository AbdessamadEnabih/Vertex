FROM golang:1.23-alpine

WORKDIR /opt/.vertex

COPY . .

RUN go build -o .bin/vertex ./cmd/

EXPOSE 64800

ENV PATH="/opt/.vertex/.bin:$PATH"

