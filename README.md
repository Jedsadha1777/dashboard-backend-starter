# Go Dashboard Backend Starter

This is a starter project for building a secure, scalable, and modular backend for admin dashboards using Go.

## Features
- JWT-based Authentication (github.com/golang-jwt/jwt/v5)
- Modular clean architecture (Controller/Service/Model)
- PostgreSQL + GORM (gorm.io/gorm)
- Middleware for auth and context
- Config via `.env` using github.com/joho/godotenv
- Seeder for initial admin
- Go Modules for dependency management

## Folder Structure
```
.
├── config/           # Load env & DB config
├── controllers/      # Handle request logic
├── db/               # Database connection + Seeder
├── middleware/       # Auth middleware (JWT)
├── models/           # GORM models
├── routes/           # Router setup
├── services/         # Business logic
├── utils/            # JWT utils, helpers
├── main.go           # Entry point
├── .env              # (not committed) Config secrets
├── go.mod / sum      # Go Modules
```

## Requirements
- Go 1.20+
- PostgreSQL (local or remote)
- Git

## Quick Start

### 1. Clone this repo
```
git clone https://github.com/your-username/go-dashboard-backend.git
cd go-dashboard-backend
```

### 2. Create `.env` file
And edit values:
```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=dashboard_db
DB_TIMEZONE=Asia/Bangkok
JWT_SECRET=a_super_secret_key
```

### 3. Install dependencies
Or manually:
```
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
```

### 4. Run the server
Option 1: Run normally
```
go run main.go
```
Server will start on `http://localhost:8080`

Option 2: Run with hot reload using air

Install air (only once):
```
go install github.com/cosmtrek/air@latest
```
Ensure that $GOPATH/bin or $HOME/go/bin is in your $PATH

Create .air.toml config (optional but recommended):
```
cp .air.example.toml .air.toml
```
Run:
```
air
```
The server will reload automatically when files change.

## API Authentication
Use Bearer Token via JWT.

### Example Login
```
curl -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{"email": "admin@example.com", "password": "password"}'
```

### Access protected route
```
curl http://localhost:8080/api/dashboard \
-H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Seeder
Seeder will automatically create the first admin if none exists in the `admins` table when running `main.go`.

## Notes
- Do not commit `.env`
- For production: change `JWT_SECRET`, use HTTPS, implement token revocation

