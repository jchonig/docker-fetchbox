FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o fetchbox .

FROM ghcr.io/linuxserver/baseimage-alpine:3.21
COPY --from=builder /app/fetchbox /usr/local/bin/fetchbox
COPY root/ /
VOLUME ["/config"]
