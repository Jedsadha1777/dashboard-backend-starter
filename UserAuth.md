# Authentication System Guide

คู่มือระบบ Authentication แบบ Clean Architecture

## 🔐 ภาพรวมระบบ

ระบบรองรับ 3 ประเภทผู้ใช้:

1. **Admin** - ผู้ดูแลระบบ
2. **User** - ผู้ใช้งานทั่วไป  
3. **Device** - อุปกรณ์ IoT

## 🏗️ โครงสร้าง Authentication

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Handler    │────▶│   Service    │────▶│  Repository  │
│              │     │              │     │              │
│ • Validate   │     │ • Business   │     │ • Database   │
│ • Response   │     │   Logic      │     │ • CRUD       │
└──────────────┘     └──────────────┘     └──────────────┘
```

### ไฟล์สำคัญ

```
internal/
├── domain/auth/
│   ├── entity/         # Admin, Token entities
│   └── repository/     # Repository interfaces
├── application/
│   └── services/
│       ├── auth_service.go    # Admin auth
│       └── user_service.go    # User auth
└── interfaces/http/
    ├── handlers/
    │   ├── auth_handler.go    # Admin endpoints
    │   └── user_handler.go    # User endpoints
    └── middleware/
        └── auth.go            # JWT validation
```

## 🔑 JWT Token Structure

### Access Token (30 นาที)
```json
{
  "user_id": 123,
  "user_type": "admin|user|device",
  "token_version": 5,
  "exp": 1725801600,
  "iat": 1725799800
}
```

### Refresh Token (1 ปี)
- เก็บใน database
- สามารถ revoke ได้
- ใช้ refresh access token

## 🛡️ Security Features

### 1. Password Security
- ความยาวขั้นต่ำ 12 ตัวอักษร (กำหนดได้)
- ต้องมีตัวพิมพ์ใหญ่, เล็ก, ตัวเลข, อักขระพิเศษ
- ไม่อนุญาต password ที่ใช้บ่อย
- ไม่อนุญาตตัวอักษรเรียงกัน (abc, 123)

### 2. Token Versioning
- แต่ละ user มี `token_version`
- เพิ่มค่าเมื่อ login/logout/เปลี่ยน password
- Token เก่าจะใช้ไม่ได้ทันที

### 3. Rate Limiting
```env
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/auth/login,/api/v1/user/auth/register
```

## 💻 Implementation Examples

### Login Flow
```go
// 1. Handler รับ request
func (h *AuthHandler) Login(c *gin.Context) {
    var input dto.LoginInput
    c.ShouldBindJSON(&input)
    
    // 2. เรียก Service
    response, err := h.authService.LoginAdmin(input)
    
    // 3. Return token
    c.JSON(200, response)
}

// 2. Service จัดการ business logic
func (s *AuthService) LoginAdmin(input dto.LoginInput) (*dto.LoginResponse, error) {
    // ค้นหา admin
    admin, err := s.adminRepo.GetByEmail(input.Email)
    
    // ตรวจสอบ password
    bcrypt.CompareHashAndPassword(admin.Password, input.Password)
    
    // Update token version
    admin.TokenVersion++
    s.adminRepo.Update(admin)
    
    // Generate tokens
    token := utils.GenerateToken(admin.ID, "admin", admin.TokenVersion)
    refreshToken := s.CreateRefreshToken(admin.ID, "admin")
    
    return &dto.LoginResponse{token, refreshToken}
}
```

### Middleware Protection
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. ดึง token
        token := c.GetHeader("Authorization")
        
        // 2. Parse และ validate
        userID, userType, tokenVer := utils.ParseToken(token)
        
        // 3. ตรวจสอบ token version
        user := getUserFromDB(userID)
        if tokenVer != user.TokenVersion {
            c.AbortWithStatusJSON(401, "Token revoked")
            return
        }
        
        // 4. Set context
        c.Set("user_id", userID)
        c.Set("user_type", userType)
        c.Next()
    }
}
```

### Role-Based Access
```go
// Admin only
func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("user_type") != "admin" {
            c.AbortWithStatusJSON(403, "Admin only")
        }
        c.Next()
    }
}

// User only
func UserRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.GetString("user_type") != "user" {
            c.AbortWithStatusJSON(403, "User only")
        }
        c.Next()
    }
}
```

## 🔄 Token Refresh

```go
func RefreshToken(refreshToken string) (*LoginResponse, error) {
    // 1. ค้นหา refresh token
    token := tokenRepo.GetByToken(refreshToken)
    
    // 2. ตรวจสอบ expiry และ revoked
    if token.IsRevoked || token.ExpiresAt.Before(time.Now()) {
        return nil, errors.New("Invalid refresh token")
    }
    
    // 3. Generate new access token
    newToken := utils.GenerateToken(token.UserID, token.UserType)
    
    return &LoginResponse{Token: newToken}
}
```

## 📝 Database Schema

### Admins Table
```sql
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    token_version INT DEFAULT 1,
    last_login TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    token_version INT DEFAULT 1,
    admin_id INT REFERENCES admins(id),
    last_login TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id INT NOT NULL,
    user_type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP
);
```

## 🚀 การใช้งาน

### 1. Admin Login
```bash
POST /api/v1/auth/login
{
  "email": "admin@example.com",
  "password": "password"
}

Response:
{
  "token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

### 2. User Registration
```bash
POST /api/v1/user/auth/register
{
  "name": "John Doe",
  "email": "john@example.com", 
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

### 3. Protected Request
```bash
GET /api/v1/admin/users
Authorization: Bearer eyJ...
```

## ⚙️ Configuration

```env
# JWT
JWT_SECRET=your_secret_key_min_32_chars
JWT_EXPIRY_MINUTES=30

# Security
SECURITY_MIN_PASSWORD_LENGTH=12

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/auth/login

# Development
ENVIRONMENT=development
AUTO_SEED=true
```

## 🔍 Troubleshooting

### Token ถูก revoke
- ตรวจสอบ token_version ใน database
- User อาจ logout หรือเปลี่ยน password

### Rate limit exceeded
- รอ 1 นาทีแล้วลองใหม่
- ตรวจสอบ IP ใน rate limiter

### Invalid token
- ตรวจสอบ JWT_SECRET ตรงกัน
- Token อาจหมดอายุ (30 นาที)

## 📚 Related Documentation

- [API Documentation](./API.md) - API endpoints reference
- [Development Guide](./DEVELOPMENT.md) - Development setup