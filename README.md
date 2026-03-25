# Modjot: Main Service

// generate mock service
mockery --name=Service --dir=internal/category/service --output=internal/category/mocks

// test in local
go test ./internal/summary/handler -v
