# ---- Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app

# Better caching
RUN go mod tidy
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main-service ./cmd/server

# ---- Minimal runtime
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /app/main-service /main-service
USER nonroot:nonroot

# Match the port your app listens on (see FIBER_PORT below)
EXPOSE 8081
ENTRYPOINT ["/main-service"]
