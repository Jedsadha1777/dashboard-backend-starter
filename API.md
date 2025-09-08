# API Documentation

Complete API endpoints reference for Dashboard Backend.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All protected endpoints require Bearer token in Authorization header:
```
Authorization: Bearer <token>
```

## Response Format
All responses follow this format:
```json
{
  "success": true,
  "data": {},
  "meta": {},
  "error": ""
}
```

---

## 🔐 Admin Authentication

### Login
- **POST** `/auth/login`
- **Body**:
```json
{
  "email": "admin@example.com",
  "password": "password"
}
```
- **Response**: Token, refresh token, expiry

### Logout
- **POST** `/auth/logout`
- **Auth**: Required (Admin)

### Refresh Token
- **POST** `/auth/refresh`
- **Body**:
```json
{
  "refresh_token": "..."
}
```

### Get Profile
- **GET** `/auth/profile`
- **Auth**: Required (Admin)

---

## 👤 User Authentication

### Register
- **POST** `/user/auth/register`
- **Body**:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

### Login
- **POST** `/user/auth/login`
- **Body**:
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

### Logout
- **POST** `/user/auth/logout`
- **Auth**: Required (User)

### Get Profile
- **GET** `/user/auth/profile`
- **Auth**: Required (User)

### Change Password
- **POST** `/user/auth/change-password`
- **Auth**: Required (User)
- **Body**:
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
- **GET** `/admin/users`
- **Auth**: Required (Admin)
- **Query Params**:
  - `page` (default: 1)
  - `limit` (default: 10, max: 100)
  - `search` (search by name/email)

### Create User
- **POST** `/admin/users`
- **Auth**: Required (Admin)
- **Body**:
```json
{
  "name": "New User",
  "email": "newuser@example.com"
}
```
- **Response**: User data + temporary password

### Get User
- **GET** `/admin/users/:id`
- **Auth**: Required (Admin)

### Update User
- **PUT** `/admin/users/:id`
- **Auth**: Required (Admin)
- **Body**:
```json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

### Delete User
- **DELETE** `/admin/users/:id`
- **Auth**: Required (Admin)

### Reset Password
- **POST** `/admin/users/:id/reset-password`
- **Auth**: Required (Admin)

---

## 📱 Device Management

### Device Authentication
- **POST** `/auth/device`
- **Body**:
```json
{
  "device_id": "DEVICE001",
  "api_key": "device_api_key_here"
}
```

### List Devices
- **GET** `/admin/devices`
- **Auth**: Required (Admin)
- **Query Params**:
  - `page` (default: 1)
  - `limit` (default: 10)
  - `search`

### Register Device
- **POST** `/admin/devices`
- **Auth**: Required (Admin)
- **Body**:
```json
{
  "device_id": "DEVICE001",
  "name": "Temperature Sensor"
}
```
- **Response**: Device data + API key

### Get Device
- **GET** `/admin/devices/:id`
- **Auth**: Required (Admin)

### Update Device
- **PUT** `/admin/devices/:id`
- **Auth**: Required (Admin)

### Delete Device
- **DELETE** `/admin/devices/:id`
- **Auth**: Required (Admin)

### Reset API Key
- **POST** `/admin/devices/:id/reset-key`
- **Auth**: Required (Admin)

---

## 📝 Article Management

### List Articles
- **GET** `/admin/articles`
- **Auth**: Required (Admin)
- **Query Params**:
  - `page` (default: 1)
  - `limit` (default: 10)
  - `search`
  - `status` (draft/published/archived)

### Create Article
- **POST** `/admin/articles`
- **Auth**: Required (Admin)
- **Body**:
```json
{
  "title": "Article Title",
  "content": "Article content...",
  "slug": "article-slug",
  "summary": "Brief summary",
  "status": "draft"
}
```

### Get Article
- **GET** `/admin/articles/:id`
- **Auth**: Required (Admin)

### Update Article
- **PUT** `/admin/articles/:id`
- **Auth**: Required (Admin)

### Delete Article
- **DELETE** `/admin/articles/:id`
- **Auth**: Required (Admin)

### Publish Article
- **POST** `/admin/articles/:id/publish`
- **Auth**: Required (Admin)

---

## 🏠 User Dashboard

### Get Dashboard
- **GET** `/user/dashboard`
- **Auth**: Required (User)

### Update Profile
- **PUT** `/user/profile`
- **Auth**: Required (User)
- **Body**:
```json
{
  "name": "Updated Name",
  "email": "newemail@example.com"
}
```

---

## 📊 Public Endpoints

### Health Check
- **GET** `/health`
- **Response**:
```json
{
  "status": "ok"
}
```

### Metrics
- **GET** `/metrics`
- **Response**: Prometheus metrics format

### Public Articles
- **GET** `/public/articles`
- **Response**: Published articles list

---

## 🔴 Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

## 📈 Rate Limiting

Default rate limits:
- Authentication endpoints: 5 requests/minute
- Other endpoints: 60 requests/minute

Rate limiting is per IP address.