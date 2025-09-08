# Development Guide

Guidelines for developing and extending the Dashboard Backend.

## 🛠️ Development Setup

### Prerequisites
```bash
# Install Go
brew install go  # macOS
# or download from https://golang.org

# Install PostgreSQL
brew install postgresql

# Install Redis (optional)
brew install redis

# Install Air for hot reload
go install github.com/air-verse/air@latest

# Install development tools
make install-tools
```

### Environment Configuration

Create `.env` file for development:
```env
# Development Settings
ENVIRONMENT=development
AUTO_SEED=true
SEED_TEST_DATA=true

# Database
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=dashboard_dev
DB_HOST=localhost
DB_PORT=5432

# Security (use strong keys in production)
JWT_SECRET=dev_secret_key_change_in_production_min_32_chars

# Server
SERVER_PORT=8080

# Enable debug logging
LOG_LEVEL=debug
```

### Database Setup
```bash
# Create development database
createdb dashboard_dev

# Run migrations and seed
go run cmd/seed/main.go

# Or use Makefile
make migrate-up
```

### Running the Application
```bash
# With hot reload (recommended)
air

# Standard run
go run main.go

# Build and run
make build
./bin/api
```

## 🏗️ Architecture Guidelines

### Clean Architecture Principles

1. **Dependency Rule**: Dependencies point inward
   - Domain → Application → Infrastructure → Interface

2. **Layer Responsibilities**:
   - **Domain**: Business entities and rules
   - **Application**: Use cases and orchestration
   - **Infrastructure**: External services (DB, cache, email)
   - **Interface**: HTTP handlers and middleware

### Adding New Features

#### 1. Create Domain Entity
```go
// internal/domain/feature/entity/feature.go
package entity

type Feature struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 2. Define Repository Interface
```go
// internal/domain/feature/repository/feature.go
package repository

type FeatureRepository interface {
    GetByID(id uint) (*entity.Feature, error)
    Create(feature *entity.Feature) error
    Update(feature *entity.Feature) error
    Delete(id uint) error
}
```

#### 3. Implement Repository
```go
// internal/infrastructure/database/feature_repository.go
package database

type featureRepository struct {
    db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) repository.FeatureRepository {
    return &featureRepository{db: db}
}
```

#### 4. Create Application Service
```go
// internal/application/services/feature_service.go
package services

type FeatureApplicationService struct {
    repo repository.FeatureRepository
}

func NewFeatureApplicationService(repo repository.FeatureRepository) *FeatureApplicationService {
    return &FeatureApplicationService{repo: repo}
}
```

#### 5. Add HTTP Handler
```go
// internal/interfaces/http/handlers/feature_handler.go
package handlers

type FeatureHandler struct {
    service *services.FeatureApplicationService
}

func NewFeatureHandler(service *services.FeatureApplicationService) *FeatureHandler {
    return &FeatureHandler{service: service}
}
```

#### 6. Update Container
```go
// internal/interfaces/http/container.go
// Add to NewContainer():
featureRepo := database.NewFeatureRepository(db)
featureService := services.NewFeatureApplicationService(featureRepo)
featureHandler := handlers.NewFeatureHandler(featureService)
```

#### 7. Add Routes
```go
// routes/routes.go
features := admin.Group("/features")
{
    features.GET("", container.FeatureHandler.List)
    features.POST("", container.FeatureHandler.Create)
    // ...
}
```

#### 8. Update Model Registry
```go
// db/model.go
modelRegistry = []interface{}{
    // existing models...
    &featureEntity.Feature{},
}
```

## 🧪 Testing

### Unit Testing
```go
// internal/application/services/feature_service_test.go
func TestFeatureService_Create(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockFeatureRepository()
    service := NewFeatureApplicationService(mockRepo)
    
    // Act
    result, err := service.Create(input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Integration Testing
```go
// tests/feature_test.go
func TestFeatureAPI(t *testing.T) {
    router := setupTestRouter()
    
    // Test Create
    req := httptest.NewRequest("POST", "/api/v1/features", body)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### Test Database
```go
func setupTestDB() *gorm.DB {
    // Use SQLite for tests
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&entity.Feature{})
    return db
}
```

## 🔍 Debugging

### Enable Debug Logging
```env
LOG_LEVEL=debug
```

### Database Query Logging
```go
// db/db.go
newLogger := logger.New(
    log.New(log.Writer(), "\r\n", log.LstdFlags),
    logger.Config{
        LogLevel: logger.Info, // Change to logger.Info for query logs
    },
)
```

### Request/Response Logging
```go
// Add to middleware
func LoggerMiddleware() gin.HandlerFunc {
    return gin.Logger()
}
```

## 📦 Dependencies Management

### Adding Dependencies
```bash
go get github.com/package/name
go mod tidy
```

### Updating Dependencies
```bash
go get -u ./...
go mod tidy
```

### Vendor Dependencies
```bash
go mod vendor
```

## 🚀 Deployment Preparation

### Build for Production
```bash
# Linux
CGO_ENABLED=0 GOOS=linux go build -o bin/api main.go

# With Makefile
make build-prod
```

### Production Checklist
- [ ] Set `ENVIRONMENT=production`
- [ ] Set `AUTO_SEED=false`
- [ ] Use strong `JWT_SECRET`
- [ ] Configure proper database credentials
- [ ] Set up SSL/TLS
- [ ] Configure rate limiting
- [ ] Set up monitoring (Prometheus)
- [ ] Configure log aggregation
- [ ] Set up backup strategy

## 🔧 Makefile Commands

```bash
make run           # Run application
make build         # Build binary
make test          # Run tests
make test-coverage # Generate coverage report
make docs          # Generate Swagger docs
make migrate-up    # Run migrations
make migrate-down  # Rollback migration
make clean         # Clean build artifacts
make fmt           # Format code
make lint          # Run linter
```

## 📝 Code Style

### Naming Conventions
- Files: `snake_case.go`
- Packages: lowercase
- Interfaces: `PascalCase` with suffix (e.g., `UserRepository`)
- Structs: `PascalCase`
- Functions/Methods: `PascalCase` (exported) or `camelCase` (private)
- Constants: `PascalCase` or `UPPER_SNAKE_CASE`

### Comments
```go
// UserService handles user-related business logic
type UserService struct {
    // ...
}

// CreateUser creates a new user with the given input
// Returns the created user and any error encountered
func (s *UserService) CreateUser(input dto.UserInput) (*entity.User, error) {
    // Validate input
    // ...
}
```

### Error Handling
```go
// Define errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrEmailExists  = errors.New("email already exists")
)

// Return wrapped errors
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

## 🐛 Common Issues

### Port Already in Use
```bash
# Find process using port
lsof -i :8080
# Kill process
kill -9 <PID>
```

### Database Connection Failed
- Check PostgreSQL is running
- Verify credentials in `.env`
- Check database exists

### Migration Failed
- Check for syntax errors
- Verify model definitions
- Review foreign key constraints

### Hot Reload Not Working
- Ensure Air is installed
- Check `.air.toml` configuration
- Verify file permissions



## 📚 Related Documentation

- [API Documentation](./API.md)** - Complete API endpoints reference
- [Adding New Features Guide](./ADDING_FEATURES.md)** - How to Add New Features: A Step-by-Step Guide
- [Authentication Guide](./UserAuth.md)** - Authentication implementation details