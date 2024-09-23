FROM golang:1.23-alpine

WORKDIR /opt/.vertex

COPY . .

RUN go build -o .bin/vertex ./cmd/

    EXPOSE 6480

ENV PATH="/opt/.vertex/.bin:$PATH"

ENTRYPOINT ["vertex","serve"]
