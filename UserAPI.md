# User API Documentation - Clean Architecture

This document outlines the user management API endpoints in the refactored Clean Architecture system.

## 📋 Architecture Overview

The User API is built using Clean Architecture with the following layers:

- **Domain Layer**: `internal/domain/user/` - User entities and business rules
- **Application Layer**: `internal/application/services/user_service.go` - Use cases
- **Infrastructure Layer**: `internal/infrastructure/database/user_repository.go` - Data access
- **Interface Layer**: `internal/interfaces/http/handlers/user_handler.go` - HTTP handlers

## 🔐 Authentication Domains

The system supports three distinct authentication domains:

1. **Admin Domain**: Administrative users (`/api/v1/auth/`)
2. **User Domain**: Regular application users (`/api/v1/user/auth/`)
3. **Device Domain**: IoT devices (`/api/v1/auth/device`)

## 🚀 User Authentication Endpoints

### User Registration

Register a new user account with strong password validation.

- **URL**: `/api/v1/user/auth/register`
- **Method**: `POST`
- **Auth Required**: No
- **Handler**: `UserHandler.Register`

**Request Body**:
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

**Success Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-09-09T10:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

**Error Response (400 Bad Request)**:
```json
{
  "success": false,
  "error": "Password not strong enough: Password must contain at least one special character"
}
```

### User Login

Authenticate an existing user.

- **URL**: `/api/v1/user/auth/login`
- **Method**: `POST`
- **Auth Required**: No
- **Handler**: `UserHandler.Login`

**Request Body**:
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-09-09T10:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

### User Logout

Invalidate user session by incrementing token version.

- **URL**: `/api/v1/user/auth/logout`
- **Method**: `POST`
- **Auth Required**: Yes (User)
- **Middleware**: `AuthMiddleware`, `UserRequired`

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "message": "User logged out successfully",
    "user_id": 1
  }
}
```

### Refresh Token

Generate new access token using valid refresh token.

- **URL**: `/api/v1/user/auth/refresh`
- **Method**: `POST`
- **Auth Required**: No (requires refresh token)
- **Handler**: `AuthHandler.RefreshToken` (shared with admin)

**Request Body**:
```json
{
  "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
}
```

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-09-09T11:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

## 👤 User Profile Management

### Get User Profile

Retrieve the authenticated user's profile.

- **URL**: `/api/v1/user/auth/profile`
- **Method**: `GET`
- **Auth Required**: Yes (User)
- **Handler**: `UserHandler.GetProfile`

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "admin_id": 0,
    "created_at": "2025-09-01T10:00:00Z",
    "updated_at": "2025-09-01T10:00:00Z",
    "last_login": "2025-09-08T10:00:00Z"
  }
}
```

### Update User Profile

Update the authenticated user's profile information.

- **URL**: `/api/v1/user/profile`
- **Method**: `PUT`
- **Auth Required**: Yes (User)
- **Handler**: `UserHandler.UpdateProfile`

**Request Body**:
```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "admin_id": 0,
    "created_at": "2025-09-01T10:00:00Z",
    "updated_at": "2025-09-08T11:00:00Z",
    "last_login": "2025-09-08T10:00:00Z"
  }
}
```

### Change Password

Change the authenticated user's password.

- **URL**: `/api/v1/user/auth/change-password`
- **Method**: `POST`
- **Auth Required**: Yes (User)
- **Handler**: `UserHandler.ChangePassword`

**Request Body**:
```json
{
  "current_password": "SecurePass123!",
  "new_password": "EvenMoreSecure456!",
  "confirm_password": "EvenMoreSecure456!"
}
```

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "message": "Password updated successfully"
  }
}
```

## 👨‍💼 Admin User Management

### List Users

Retrieve paginated list of users (Admin only).

- **URL**: `/api/v1/admin/users`
- **Method**: `GET`
- **Auth Required**: Yes (Admin)
- **Handler**: `UserHandler.ListUsers`

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `search`: Search term for name or email

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "admin_id": 1,
      "created_at": "2025-09-01T10:00:00Z",
      "updated_at": "2025-09-01T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

### Create User (Admin)

Create a new user account by admin with temporary password.

- **URL**: `/api/v1/admin/users`
- **Method**: `POST`
- **Auth Required**: Yes (Admin)
- **Handler**: `UserHandler.CreateUser`

**Request Body**:
```json
{
  "name": "New User",
  "email": "new.user@example.com"
}
```

**Success Response (201 Created)**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 3,
      "name": "New User",
      "email": "new.user@example.com",
      "admin_id": 1,
      "created_at": "2025-09-08T12:00:00Z",
      "updated_at": "2025-09-08T12:00:00Z"
    },
    "temporary_password": "Rand0mP@ssw0rd",
    "message": "User created successfully. Please inform the user to change their password after first login."
  }
}
```

### Get User by ID

Retrieve specific user information (Admin only).

- **URL**: `/api/v1/admin/users/:id`
- **Method**: `GET`
- **Auth Required**: Yes (Admin)
- **Handler**: `UserHandler.GetUser`

**Success Response (200 OK)**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "admin_id": 1,
    "created_at": "2025-09-01T10:00:00Z",
    "updated_at": "2025-09-01T10:00:00Z"
  }
}
```

## 🔒 Security Features

### Password Strength Validation

The system enforces strong password requirements:

- Minimum 12 characters (configurable)
- At least one lowercase letter
- At least one uppercase letter
- At least one digit
- At least one special character
- No common passwords (password, 123456, etc.)
- No sequential characters (abc, 123, etc.)

### Token Management

- **Access Token**: 30 minutes validity (configurable)
- **Refresh Token**: 1 year validity
- **Token Versioning**: Increment on login/logout to invalidate old tokens
- **Secure Storage**: Refresh tokens stored in database with revocation support

### Rate Limiting

Configurable rate limiting on authentication endpoints:

```env
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/user/auth/login,/api/v1/user/auth/register
```

## 🏗️ Architecture Components

### Domain Layer
```
internal/domain/user/
├── entity/user.go          # User domain entity
├── repository/user.go      # Repository interface
├── service/                # Domain services (if needed)
└── errors/                 # Domain-specific errors
```

### Application Layer
```
internal/application/
├── dto/user.go            # User DTOs
├── services/user_service.go # User use cases
└── ports/                 # External service interfaces
```

### Infrastructure Layer
```
internal/infrastructure/
└── database/user_repository.go # Repository implementation
```

### Interface Layer
```
internal/interfaces/http/
├── handlers/user_handler.go    # HTTP handlers
├── middleware/user.go          # User-specific middleware
└── dto/                       # HTTP-specific DTOs
```

## 📝 Error Responses

### Authentication Errors (401 Unauthorized)
```json
{
  "success": false,
  "error": "Invalid email or password"
}
```

### Validation Errors (400 Bad Request)
```json
{
  "success": false,
  "error": "Name is required; Email must be a valid email address"
}
```

### Permission Errors (403 Forbidden)
```json
{
  "success": false,
  "error": "Forbidden: you don't have permission to access this resource"
}
```

### Rate Limit Errors (429 Too Many Requests)
```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later"
}
```

## 🧪 Testing

### Unit Tests
```bash
# Test user domain
go test ./internal/domain/user/... -v

# Test user application services
go test ./internal/application/services/... -v

# Test user handlers
go test ./internal/interfaces/http/handlers/... -v
```

### Integration Tests
```bash
# Test user API endpoints
go test ./tests/user_api_test.go -v
```

## 🔄 User Types

### Self-Registered Users
- Register through `/api/v1/user/auth/register`
- `admin_id` is 0 or NULL
- Set their own passwords
- Full control over their profiles

### Admin-Created Users
- Created by admins through `/api/v1/admin/users`
- Have `admin_id` referencing the creating admin
- Receive temporary passwords
- Can only be managed by the creating admin

## 📚 Related Documentation

- [Authentication System](UserAuth.md) - Detailed authentication flow
- [Clean Architecture Guide](README.md) - Overall system architecture
- [API Reference](swagger/index.html) - Interactive API documentation