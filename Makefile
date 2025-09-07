.PHONY: run build test docs clean install-tools docker-up docker-down

# Development commands
run:
	go run main.go

build:
	go build -o bin/api main.go

test:
	go test ./tests/... -v

test-coverage:
	go test ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Documentation
docs:
	swag init
	@echo "Swagger docs generated at docs/"

# Database
migrate-up:
	go run cmd/migrate/main.go -direction=up

migrate-down:
	go run cmd/migrate/main.go -direction=down -steps=1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Development tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Cleanup
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -modcache

# Linting and formatting
fmt:
	go fmt ./...

lint:
	golangci-lint run

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api main.go

# Quick development setup
dev-setup: install-tools
	cp .env.example .env
	@echo "Please edit .env file with your configuration"
	@echo "Then run: make docker-up && make run"