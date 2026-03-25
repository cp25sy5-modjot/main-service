# Modjot: Main Service

// generate mock service
mockery --name=Service --dir=internal/user/service --output=internal/user/mocks

// test in local
go test ./internal/user/handler -v
