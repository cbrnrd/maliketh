# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/admin ./cmd/admin
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/c2 ./cmd/c2
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/createop ./cmd/createop

# Admin server image
FROM alpine:3.19 AS admin

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/bin/admin /app/admin
COPY --from=builder /build/bin/createop /app/createop
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/configs /app/configs

EXPOSE 5000

ENTRYPOINT ["/app/admin"]

# C2 server image
FROM alpine:3.19 AS c2

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/bin/c2 /app/c2
COPY --from=builder /build/configs /app/configs

EXPOSE 8080

ENTRYPOINT ["/app/c2"]

# Migration tool image
FROM alpine:3.19 AS migrate

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/bin/migrate /app/migrate
COPY --from=builder /build/migrations /app/migrations

ENTRYPOINT ["/app/migrate"]

# Default target: all-in-one image with all binaries
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata bash

WORKDIR /app

COPY --from=builder /build/bin/admin /app/admin
COPY --from=builder /build/bin/c2 /app/c2
COPY --from=builder /build/bin/migrate /app/migrate
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/configs /app/configs

EXPOSE 5000 8080

# Default to running the admin server
CMD ["/app/admin"]
