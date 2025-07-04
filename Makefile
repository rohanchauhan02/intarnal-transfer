# Makefile for the Internal Transfer Project

.PHONY: postgres setup start build test clean

# Start a local PostgreSQL instance using Docker
postgres:
	docker pull postgres:latest
	@if docker ps -a --format '{{.Names}}' | grep -Eq '^postgres$$'; then docker rm -f postgres; fi
	docker run --name postgres \
		-e POSTGRES_USER=root \
		-e POSTGRES_PASSWORD=root \
		-e POSTGRES_DB=internal_transfer_local \
		-p 5432:5432 \
		-d postgres:latest

# Set up Go modules and configuration files, and start PostgreSQL
setup: postgres
	go mod tidy
	cp configs/app.config.sample.yml configs/app.config.local.yml

# Run the application
start:
	go run app/main.go

# Build the application binary
build:
	go build -o app/main app/main.go

# Run tests with coverage reporting
test:
	go test -v ./domain/banking/usecase -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts and coverage files
clean:
	rm -rf app/main coverage.out coverage.html

install-mockgen:
	go install github.com/golang/mock/mockgen@latest

mock:
	mockgen -source=domain/banking/banking.go -destination=file/mocks/mock_banking/usecase.go

