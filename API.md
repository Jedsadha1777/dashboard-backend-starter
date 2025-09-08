# API Documentation

Base URL: `http://localhost:3000/api/v1`

## Authentication

All protected endpoints require Bearer token in Authorization header:
```
Authorization: Bearer <token>
```

## Response Format

```json
{
  "success": true,
  "data": {},
  "error": "",
  "meta": {}
}
```

---

## 🔐 Admin Authentication

### Login
**POST** `/auth/login`
```json
{
  "email": "admin@example.com",
  "password": "password"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_at": "2024-01-01T12:00:00Z",
    "user_id": 1,
    "user_type": "admin"
  }
}
```

### Refresh Token
**POST** `/auth/refresh`
```json
{
  "refresh_token": "eyJ..."
}
```

### Logout
**POST** `/auth/logout` (Protected)

### Get Profile
**GET** `/auth/profile` (Protected)

---

## 👤 User Authentication

### Register
**POST** `/user/auth/register`
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

### Login
**POST** `/user/auth/login`
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

### Change Password
**POST** `/user/auth/change-password` (Protected)
```json
{
  "current_password": "OldPass123!",
  "new_password": "NewPass456!",
  "confirm_password": "NewPass456!"
}
```

---

## 👥 User Management (Admin Only)

### List Users
**GET** `/admin/users`

Query Parameters:
- `page` (default: 1)
- `limit` (default: 10, max: 100)
- `search` - Search by name or email

### Create User
**POST** `/admin/users`
```json
{
  "name": "New User",
  "email": "newuser@example.com"
}
```

**Response includes temporary password**

### Get User
**GET** `/admin/users/:id`

### Update User
**PUT** `/admin/users/:id`
```json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

### Delete User
**DELETE** `/admin/users/:id`

---

## 📱 Device Management (Admin Only)

### Register Device
**POST** `/admin/devices`
```json
{
  "device_id": "DEVICE001",
  "name": "Temperature Sensor"
}
```

**Response includes API key**

### List Devices
**GET** `/admin/devices`

Query Parameters:
- `page` (default: 1)
- `limit` (default: 10)
- `search`

### Get Device
**GET** `/admin/devices/:id`

### Reset API Key
**POST** `/admin/devices/:id/reset-key`

### Device Authentication
**POST** `/auth/device`
```json
{
  "device_id": "DEVICE001",
  "api_key": "generated_api_key"
}
```

---

## 📝 Article Management (Admin Only)

### Create Article
**POST** `/admin/articles`
```json
{
  "title": "Article Title",
  "content": "Article content...",
  "slug": "article-slug",
  "summary": "Brief summary",
  "status": "draft"
}
```

### List Articles
**GET** `/admin/articles`

Query Parameters:
- `page` (default: 1)
- `limit` (default: 10)
- `search`
- `status` (draft/published/archived)

### Get Article
**GET** `/admin/articles/:id`

### Update Article
**PUT** `/admin/articles/:id`

### Delete Article
**DELETE** `/admin/articles/:id`

### Publish Article
**POST** `/admin/articles/:id/publish`

---

## 🏠 User Dashboard

### Get Dashboard
**GET** `/user/dashboard` (User Auth Required)

### Update Profile
**PUT** `/user/profile` (User Auth Required)
```json
{
  "name": "Updated Name",
  "email": "newemail@example.com"
}
```

---

## 📊 Public Endpoints

### Health Check
**GET** `/health`
```json
{
  "status": "healthy",
  "database": "connected",
  "uptime": 3600
}
```

### Readiness Check
**GET** `/health/ready`
```json
{
  "status": "ready",
  "checks": {
    "database": {
      "status": "healthy",
      "pool": {
        "open": 10,
        "in_use": 2,
        "idle": 8
      }
    },
    "memory": {
      "alloc_mb": 12,
      "goroutines": 15
    }
  }
}
```

### Public Articles
**GET** `/public/articles`

---

## 🔴 Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found |
| 409 | Conflict - Resource already exists |
| 429 | Too Many Requests |
| 500 | Internal Server Error |

## 📈 Rate Limiting

- Authentication endpoints: 5 requests/minute
- Other endpoints: 60 requests/minute

Rate limiting is per IP address.