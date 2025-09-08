# Go Dashboard Backend Starter - Clean Architecture

โปรเจกต์ backend สำหรับ Admin Dashboard ที่ใช้ Clean Architecture และ Domain-Driven Design

## 🏗️ Architecture Overview

โปรเจกต์นี้ใช้ **Clean Architecture** พร้อม **Domain-Driven Design** เพื่อความยืดหยุ่น ทดสอบง่าย และบำรุงรักษาได้

```
dashboard-starter/
├── cmd/                          # Application entry points
│   ├── api/main.go              # HTTP API server
│   └── seed/main.go             # Database seeding
├── internal/                     # Private application code
│   ├── domain/                   # Domain layer (business logic)
│   │   ├── auth/                # Authentication domain
│   │   ├── user/                # User management domain
│   │   ├── device/              # IoT device domain
│   │   ├── article/             # Content management domain
│   │   └── shared/              # Shared domain models
│   ├── application/              # Application layer (use cases)
│   │   ├── dto/                 # Data transfer objects
│   │   ├── ports/               # Interface definitions
│   │   └── services/            # Application services
│   ├── infrastructure/           # Infrastructure layer
│   │   ├── database/            # Database implementations
│   │   ├── cache/               # Cache implementations
│   │   ├── email/               # Email service
│   │   └── logger/              # Logging service
│   └── interfaces/               # Interface adapters
│       └── http/                # HTTP layer
│           ├── handlers/        # HTTP handlers
│           └── middleware/      # HTTP middleware
├── pkg/                         # Public packages
├── configs/                     # Configuration files
├── migrations/                  # Database migrations
└── tests/                       # Integration tests
```

## ✨ คุณสมบัติหลัก

### **Authentication & Authorization**
- **3 Authentication Domains**: Admin, User, Device
- JWT + Refresh Token พร้อม Token Versioning
- Role-based access control
- Strong password validation
- API Key สำหรับ IoT devices

### **Clean Architecture Benefits**
- **Domain-Driven Design**: แยก business domains ชัดเจน
- **Dependency Inversion**: Infrastructure ไม่ขึ้นกับ Domain
- **Testability**: แต่ละ layer test ได้อิสระ
- **Maintainability**: เพิ่มฟีเจอร์ใหม่ง่าย
- **Scalability**: พร้อมแยกเป็น microservices

### **Database & Repository**
- PostgreSQL + GORM
- Repository pattern พร้อม interfaces
- Transaction management
- Auto migrations และ seeding
- Connection pooling

### **Security & Performance**
- IP-based rate limiting
- CORS middleware
- SQL injection protection
- Trusted proxy configuration
- Password hashing with bcrypt

### **Infrastructure**
- Structured logging (Zap)
- Redis caching (optional)
- Email service (SMTP/Mock)
- Prometheus metrics
- Graceful shutdown

### **Development Experience**
- Hot reload with Air
- Swagger documentation
- Comprehensive error handling
- Validation with custom messages
- Environment-based configuration

## 🚀 Quick Start

### 1. Requirements
- Go 1.20+
- PostgreSQL 13+
- Redis (optional)

### 2. Installation

```bash
# Clone repository
git clone <repository-url>
cd dashboard-starter

# Install dependencies
go mod tidy

# Setup environment
cp .env.example .env
# Edit .env with your configuration
```

### 3. Configuration

Key environment variables:

```env
# Database
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=dashboard
DB_HOST=localhost
DB_PORT=5432

# Security
JWT_SECRET=your_strong_secret_key
SECURITY_MIN_PASSWORD_LENGTH=12

# Server
SERVER_PORT=8080
TRUSTED_PROXIES=

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/auth/login,/api/v1/user/auth/register
```

### 4. Database Setup

```bash
# Create database
createdb dashboard

# Run migrations and seed data
go run cmd/seed/main.go
```

### 5. Run Application

```bash
# Development (with hot reload)
air

# Or standard run
go run main.go
```

Server starts at `http://localhost:8080`
Swagger UI: `http://localhost:8080/swagger/index.html`

## 📋 API Endpoints

### **Authentication**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Admin login |
| POST | `/api/v1/auth/logout` | Admin logout |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/profile` | Get admin profile |

### **User Management**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/user/auth/register` | User registration |
| POST | `/api/v1/user/auth/login` | User login |
| GET | `/api/v1/user/profile` | Get user profile |
| PUT | `/api/v1/user/profile` | Update user profile |

### **Admin User Management**
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | List users (paginated) |
| POST | `/api/v1/admin/users` | Create user |
| GET | `/api/v1/admin/users/:id` | Get user by ID |
| PUT | `/api/v1/admin/users/:id` | Update user |
| DELETE | `/api/v1/admin/users/:id` | Delete user |

### **Device Management**
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/device` | Device authentication |
| GET | `/api/v1/admin/devices` | List devices |
| POST | `/api/v1/admin/devices` | Register device |
| POST | `/api/v1/admin/devices/:id/reset-key` | Reset API key |

### **Content Management**
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/articles` | List articles |
| POST | `/api/v1/admin/articles` | Create article |
| PUT | `/api/v1/admin/articles/:id` | Update article |
| POST | `/api/v1/admin/articles/:id/publish` | Publish article |

## 🧪 Testing

```bash
# Run tests
go test ./tests/... -v

# Test coverage
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 🏗️ Development

### Adding New Domain

1. **Create domain structure**:
```bash
mkdir -p internal/domain/newdomain/{entity,repository,service,errors}
```

2. **Define entities** in `internal/domain/newdomain/entity/`

3. **Create repository interface** in `internal/domain/newdomain/repository/`

4. **Implement repository** in `internal/infrastructure/database/`

5. **Create application service** in `internal/application/services/`

6. **Add HTTP handlers** in `internal/interfaces/http/handlers/`

7. **Update dependency injection** in `internal/interfaces/http/container.go`

### Code Organization Principles

- **Domain layer**: Pure business logic, no external dependencies
- **Application layer**: Use cases, orchestrates domain operations
- **Infrastructure layer**: External concerns (database, cache, email)
- **Interface layer**: HTTP handlers, middleware, DTOs

### Best Practices

- Use dependency injection
- Write tests for each layer
- Keep domain logic pure
- Use interfaces for external dependencies
- Follow SOLID principles

## 📊 Monitoring & Observability

- **Metrics**: Prometheus metrics at `/metrics`
- **Logging**: Structured logging with Zap
- **Health Check**: Available at `/health`
- **API Documentation**: Swagger UI at `/swagger/index.html`

## 🔧 Production Deployment

1. **Environment Setup**:
   - Set `GIN_MODE=release`
   - Configure strong `JWT_SECRET`
   - Set up proper `TRUSTED_PROXIES`
   - Configure Redis for caching
   - Set up SMTP for emails

2. **Database**:
   - Use connection pooling
   - Set up database backups
   - Monitor query performance

3. **Security**:
   - Use HTTPS
   - Configure rate limiting
   - Set up proper CORS
   - Monitor for suspicious activity

## 📚 Architecture Documentation

For detailed architecture information, see:
- `UserAPI.md` - User management API documentation
- `UserAuth.md` - Authentication system documentation

## 🤝 Contributing

1. Follow Clean Architecture principles
2. Write tests for new features
3. Update documentation
4. Use conventional commit messages
5. Ensure all tests pass

## 📄 License

MIT License 