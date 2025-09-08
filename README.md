# Go Dashboard Backend - Clean Architecture

Production-ready admin dashboard backend with JWT authentication, built using Clean Architecture and Domain-Driven Design.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)
- PostgreSQL 15+ (if not using Docker)

### Run with Docker (Recommended)

```bash
# Clone repository
git clone <repository-url>
cd dashboard-backend-starter

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api
```

The API will be available at:
- API: http://localhost:3000
- Swagger UI: http://localhost:3000/swagger/index.html
- Health Check: http://localhost:3000/health

### Run Locally

```bash
# Setup environment
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod tidy

# Run migrations and seed
go run cmd/seed/main.go

# Start server
go run main.go
# Or with hot reload
air
```

## 📁 Project Structure

```
dashboard-backend-starter/
├── cmd/seed/                 # Database seeding
├── config/                   # Configuration management
├── db/                       # Database initialization
├── internal/
│   ├── domain/              # Business entities & rules
│   │   ├── auth/           # Admin authentication
│   │   ├── user/           # User management
│   │   ├── device/         # Device/IoT management
│   │   └── article/        # Content management
│   ├── application/         # Use cases & services
│   ├── infrastructure/      # External services (DB, cache, email)
│   └── interfaces/http/     # HTTP handlers & middleware
├── routes/                  # API route definitions
├── utils/                   # Shared utilities
├── docker-compose.yml       # Docker orchestration
├── Dockerfile              # Container build
└── main.go                 # Application entry point
```

## 🔑 Default Credentials

After running `docker-compose up`, a default admin account is created:

```
Email: admin@example.com
Password: (check logs for generated password)
```

**Important:** Change the default password immediately after first login.

## 🛠️ Development

### Database Operations

```bash
# Run migrations
make migrate-up

# Rollback migration
make migrate-down

# Seed database
make seed

# Backup database
docker-compose exec postgres pg_dump -U postgres dashboard > backup.sql
```

### Docker Commands

```bash
# Build and start
make docker-up

# Stop all services
make docker-down

# Rebuild API
docker-compose build api

# View logs
make docker-logs

# Clean everything
make docker-clean
```

### Testing

```bash
# Run tests
go test ./...

# With coverage
go test ./... -cover

# Integration tests
make test
```

## 📚 API Documentation

Full API documentation is available at `/swagger` when the server is running.

### Key Endpoints

#### Authentication
- `POST /api/v1/auth/login` - Admin login
- `POST /api/v1/user/auth/register` - User registration
- `POST /api/v1/user/auth/login` - User login

#### Admin Operations
- `GET /api/v1/admin/users` - List users
- `POST /api/v1/admin/users` - Create user
- `GET /api/v1/admin/devices` - List devices
- `POST /api/v1/admin/articles` - Create article

## 📚 Documentation

- [API.md](./API.md) - Complete API endpoint reference
- [ADDING_FEATURES.md](./ADDING_FEATURES.md) - Step-by-step guide to add new features
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development setup and guidelines
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Production deployment guide

## 🚢 Deployment

### Production Deployment

1. **Update environment variables**
```bash
# Edit .env for production
ENVIRONMENT=production
AUTO_SEED=false
JWT_SECRET=<strong-secret-key>
```

2. **Build and deploy**
```bash
# Build production image
docker build -t dashboard-api:production .

# Deploy with Docker Compose
docker-compose -f docker-compose.yml up -d
```

3. **Setup HTTPS (optional but recommended)**
- Add Nginx reverse proxy
- Configure SSL certificates
- See [DEPLOYMENT.md](./DEPLOYMENT.md) for details

### Health Monitoring

The application provides health check endpoints:
- `/health` - Basic health status
- `/health/ready` - Detailed readiness check
- `/health/live` - Liveness probe

## 🔒 Security Features

- JWT authentication with refresh tokens
- Token versioning for instant revocation
- Password strength validation
- Rate limiting on sensitive endpoints
- SQL injection protection
- Request size limits
- Security headers (CORS, XSS, CSRF protection)

## 📄 License

MIT License