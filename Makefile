.PHONY: run build test docs clean install-tools docker-up docker-down
.PHONY: dev-migrate migrate-install migrate-up migrate-down migrate-drop migrate-create migrate-status seed
.PHONY: test-coverage fmt lint build-prod dev-setup docker-logs

# Load environment variables
include .env
export

# Database URL for migrations
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

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

# Database - Development (GORM AutoMigrate)
dev-migrate:
	go run cmd/seed/main.go

# Database - Production (golang-migrate)
migrate-install:
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	@echo "Running migrations..."
	@migrate -database "$(DB_URL)" -path migrations up

migrate-down:
	@echo "Rolling back last migration..."
	@migrate -database "$(DB_URL)" -path migrations down 1

migrate-drop:
	@echo "Dropping all migrations..."
	@migrate -database "$(DB_URL)" -path migrations drop -f

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-status:
	@echo "Current migration version:"
	@migrate -database "$(DB_URL)" -path migrations version

seed:
	go run cmd/seed/main.go

# Development tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest

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


# Docker commands
docker-build:
	docker build -t dashboard-api:latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-restart:
	docker-compose restart api

docker-clean:
	docker-compose down -v
	docker system prune -f

# Production deployment
deploy-prod:
	docker-compose --profile production up -d

# Development with monitoring
dev-with-monitoring:
	docker-compose --profile monitoring up -d

# Database backup
backup-db:
	@mkdir -p backups
	docker-compose exec postgres pg_dump -U postgres dashboard | gzip > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql.gz
	@echo "✅ Backup completed"

# Health check
health-check:
	@curl -f http://localhost:3000/health || exit 1
	@echo "✅ API is healthy"

.PHONY: docker-build docker-up docker-down docker-logs docker-restart docker-clean deploy-prod dev-with-monitoring backup-db health-check