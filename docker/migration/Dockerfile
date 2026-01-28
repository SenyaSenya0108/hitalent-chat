FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY migrations/. ./migrations

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.21

ARG UID
ARG GID

RUN addgroup -g $GID develop && adduser -u $UID -G develop -s /bin/sh -D develop

USER develop

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose