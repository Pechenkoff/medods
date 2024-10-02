FROM golang:1.22.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update

COPY --from=builder /app/server /app/server

CMD ["/app/server", "-config=./config/config.yaml", "-migration=file://./db/migrations"]