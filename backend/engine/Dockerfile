FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -ldflags "-s -w" -o engine .

FROM alpine:3.21.3
WORKDIR /app
COPY --from=builder /app/engine ./
ENTRYPOINT ["/app/engine"]