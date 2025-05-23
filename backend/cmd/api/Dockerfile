FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
WORKDIR /app/cmd/api
RUN go build -ldflags "-s -w" -o /app/api

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /app/api ./
ENTRYPOINT ["/app/api"]