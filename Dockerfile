# ---- Build stage
FROM golang:1.22-alpine AS build
WORKDIR /app

# Better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build API and Worker binaries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker

# ---- Minimal runtime
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy both binaries from build stage
COPY --from=build /app/api /api
COPY --from=build /app/worker /worker

USER nonroot:nonroot

# Match the port your Fiber app listens on (same as before)
EXPOSE 8081

# Default entrypoint = API server
ENTRYPOINT ["/api"]
