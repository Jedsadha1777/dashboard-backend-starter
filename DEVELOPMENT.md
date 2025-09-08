# Development Guide

## Setup Development Environment

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)
- Redis (optional, or use Docker)
- Air (for hot reload)

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd dashboard-backend-starter

# Install Go dependencies
go mod tidy

# Install development tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Setup environment
cp .env.example .env
# Edit .env with your local settings

# Start dependencies with Docker
docker-compose up -d postgres redis

# Run migrations and seed
go run cmd/seed/main.go

# Start with hot reload
air
```

---

## Architecture Overview

### Clean Architecture Layers

```
Interface Layer (HTTP) → Application Layer → Domain Layer → Infrastructure Layer
```

- **Interface**: HTTP handlers, middleware, request/response
- **Application**: Business logic, use cases, services
- **Domain**: Entities, business rules, repository interfaces
- **Infrastructure**: Database, external services, implementations

### Adding New Features

#### 1. Create Domain Entity
```go
// internal/domain/product/entity/product.go
type Product struct {
    ID    uint      `json:"id" gorm:"primaryKey"`
    Name  string    `json:"name" gorm:"size:255;not null"`
    Price float64   `json:"price" gorm:"type:decimal(10,2)"`
}
```

#### 2. Define Repository Interface
```go
// internal/domain/product/repository/product.go
type ProductRepository interface {
    GetByID(id uint) (*entity.Product, error)
    Create(product *entity.Product) error
    Update(product *entity.Product) error
    Delete(id uint) error
}
```

#### 3. Implement Repository
```go
// internal/infrastructure/database/product_repository.go
type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}
```

#### 4. Create Service
```go
// internal/application/services/product_service.go
type ProductService struct {
    repo repository.ProductRepository
}

func (s *ProductService) CreateProduct(input dto.ProductInput) (*entity.Product, error) {
    // Business logic here
}
```

#### 5. Add HTTP Handler
```go
// internal/interfaces/http/handlers/product_handler.go
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // Handle HTTP request/response
}
```

#### 6. Register Routes
```go
// routes/routes.go
products := admin.Group("/products")
{
    products.POST("", container.ProductHandler.CreateProduct)
    products.GET("/:id", container.ProductHandler.GetProduct)
}
```

#### 7. Register Model for Migration
```go
// db/model.go
modelRegistry = []interface{}{
    // existing models...
    &productEntity.Product{},
}
```

---

## Common Tasks

### Database Operations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_products_table

# Run migrations
make migrate-up

# Rollback
make migrate-down

# Access database
docker-compose exec postgres psql -U postgres -d dashboard
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/application/services/...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Generate Swagger Documentation

```bash
# Generate/update swagger docs
swag init

# Format with descriptions
swag init --parseDependency --parseInternal
```

### Working with Docker

```bash
# Build only API
docker-compose build api

# Run specific service
docker-compose up postgres

# Execute commands in container
docker-compose exec api sh

# View real-time logs
docker-compose logs -f api

# Clean rebuild
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

---

## Code Standards

### Project Structure
```
internal/
├── domain/           # Business logic (no external dependencies)
│   └── {feature}/
│       ├── entity/   # Business entities
│       └── repository/ # Repository interfaces
├── application/      # Use cases
│   ├── dto/         # Data transfer objects
│   └── services/    # Application services
├── infrastructure/   # External services
│   ├── database/    # Repository implementations
│   ├── cache/       # Redis implementation
│   └── email/       # Email service
└── interfaces/      # External interfaces
    └── http/        # HTTP layer
        ├── handlers/  # HTTP handlers
        └── middleware/ # HTTP middleware
```

### Naming Conventions
- Files: `snake_case.go`
- Packages: lowercase
- Exported functions/types: `PascalCase`
- Private functions/types: `camelCase`

### Error Handling
```go
// Domain errors
var ErrProductNotFound = errors.New("product not found")

// Return wrapped errors
if err != nil {
    return nil, fmt.Errorf("failed to create product: %w", err)
}

// HTTP error response
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "error": err.Error(),
    })
    return
}
```

### Validation
```go
// Use struct tags for validation
type ProductInput struct {
    Name  string  `json:"name" binding:"required" validate:"required,min=3,max=255"`
    Price float64 `json:"price" binding:"required" validate:"required,min=0"`
}

// Validate in handler
if err := c.ShouldBindJSON(&input); err != nil {
    // Handle validation error
}
```

---

## Debugging

### Enable Debug Logging
```env
LOG_LEVEL=debug
GIN_MODE=debug
```

### Database Query Logging
```go
// In db/db.go, change LogLevel
logger.Config{
    LogLevel: logger.Info, // Shows all SQL queries
}
```

### Check Running Services
```bash
# Check if services are healthy
docker-compose ps

# Check port usage
lsof -i :3000

# Check database connection
docker-compose exec postgres pg_isready
```

### Common Issues

#### Port Already in Use
```bash
# Find and kill process
lsof -i :3000
kill -9 <PID>
```

#### Database Connection Failed
```bash
# Check PostgreSQL is running
docker-compose up -d postgres

# Check credentials in .env
echo $DB_PASSWORD

# Test connection
psql -h localhost -U postgres -d dashboard
```

#### Module Issues
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy and verify
go mod tidy
go mod verify
```

---

## Git Workflow

### Branch Strategy
```bash
main          # Production-ready code
├── develop   # Integration branch
└── feature/* # Feature branches
```

### Commit Messages
```bash
feat: Add product management
fix: Resolve database connection issue
docs: Update API documentation
refactor: Simplify auth middleware
test: Add user service tests
```

### Pre-commit Checklist
- [ ] Tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Swagger updated (`swag init`)
- [ ] No linting errors (`golangci-lint run`)
- [ ] Documentation updated