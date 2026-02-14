# ---- Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project
COPY . .

# Build both API and Worker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker

# ---- Runtime stage (Alpine, NOT distroless)
FROM alpine:3.20

RUN apk add --no-cache tzdata

# Create non-root user and uploads dir
RUN adduser -D -g '' appuser && \
    mkdir -p /uploads && \
    chown -R appuser:appuser /uploads

WORKDIR /

# Copy binaries from build stage
COPY --from=build /app/api /api
COPY --from=build /app/worker /worker

# Run as unprivileged user
USER appuser

# API port
EXPOSE 8081

# Default: run API
CMD ["/api"]
