FROM golang:1.23-alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

WORKDIR /

COPY ./migrations ./migrations

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

WORKDIR /app/cmd/web

RUN CGO_ENABLED=0 GOOS=linux go build -o /api

CMD ["/api"]