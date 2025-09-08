# Go Dashboard Backend Starter - Clean Architecture

โปรเจกต์ backend สำหรับ Admin Dashboard ที่ใช้ Clean Architecture และ Domain-Driven Design

## 🏗️ Architecture Overview

โปรเจกต์นี้ใช้ **Clean Architecture** พร้อม **Domain-Driven Design** เพื่อความยืดหยุ่น ทดสอบง่าย และบำรุงรักษาได้

### Layer Architecture
```
┌─────────────────────────────────────────────┐
│          Interface Layer (HTTP)             │
│         handlers / middleware               │
├─────────────────────────────────────────────┤
│         Application Layer                   │
│      services / use cases / DTOs            │
├─────────────────────────────────────────────┤
│           Domain Layer                      │
│      entities / repositories                │
├─────────────────────────────────────────────┤
│        Infrastructure Layer                 │
│   database / cache / email / logger         │
└─────────────────────────────────────────────┘
```

## ✨ คุณสมบัติหลัก

- **Clean Architecture** - แยก business logic จาก infrastructure
- **3 Authentication Domains** - Admin, User, Device
- **JWT + Refresh Token** - พร้อม token versioning
- **Strong Security** - Rate limiting, CORS, SQL injection protection
- **PostgreSQL + GORM** - Repository pattern
- **Redis Cache** - Optional caching layer
- **Swagger Documentation** - Auto-generated API docs
- **Prometheus Metrics** - Monitoring ready
- **Hot Reload** - Development with Air

## 🚀 Quick Start

### Requirements
- Go 1.20+ 
- PostgreSQL 13+
- Redis (optional)

### Installation

```bash
# Clone repository
git clone <repository-url>
cd dashboard-starter

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Install dependencies
go mod tidy

# Create database
createdb dashboard

# Run migrations and seed
go run cmd/seed/main.go

# Start server
go run main.go
# Or with hot reload
air
```

Server starts at `http://localhost:8080`

### Important URLs
- **API Base**: `http://localhost:8080/api/v1`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:8080/metrics`

## 📁 Project Structure

```
dashboard-starter/
├── cmd/seed/             # Database seeding command
├── config/               # Configuration management
├── db/                   # Database initialization
├── docs/                 # Swagger documentation
├── internal/             
│   ├── domain/           # Business logic
│   │   ├── auth/        
│   │   ├── user/        
│   │   ├── device/      
│   │   └── article/     
│   ├── application/      # Use cases & DTOs
│   ├── infrastructure/   # External services
│   └── interfaces/http/  # HTTP handlers
├── pkg/                  # Public packages
├── routes/               # Route definitions
├── tests/                # Integration tests
├── utils/                # Utility functions
└── main.go              # Entry point
```

## 📚 Documentation

- **[API Documentation](./API.md)** - Complete API endpoints reference
- **[Authentication Guide](./UserAuth.md)** - Authentication implementation details
- **[Development Guide](./DEVELOPMENT.md)** - Development setup and guidelines

## 🧪 Testing

```bash
# Run all tests
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Using Makefile
make test
make test-coverage
```

## 🔧 Configuration

Key environment variables:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dashboard

# Security
JWT_SECRET=your_secret_key_min_32_chars
SECURITY_MIN_PASSWORD_LENGTH=12

# Development
ENVIRONMENT=development
AUTO_SEED=true
```

See `.env.example` for complete configuration options.

## 🤝 Contributing

1. Follow Clean Architecture principles
2. Write tests for new features
3. Update documentation
4. Use conventional commits

## 📄 License

MIT License