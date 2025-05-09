# User API Documentation

This document outlines the API endpoints for user registration, authentication, and management in the Dashboard Backend Starter application.

## Two Authentication Domains

The system has two distinct authentication domains:

1. **Admin Domain**: Administrative users who manage the dashboard (/api/v1/auth/...)
2. **User Domain**: Regular users of the application (/api/v1/user/auth/...)

## User Registration and Authentication Endpoints

### User Registration

Allows new users to create an account.

- **URL**: `/api/v1/user/auth/register`
- **Method**: `POST`
- **Auth Required**: No

**Request Body**:

```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!"
}
```

**Response (201 Created)**:

```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-05-09T10:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

### User Login

Authenticates a registered user.

- **URL**: `/api/v1/user/auth/login`
- **Method**: `POST`
- **Auth Required**: No

**Request Body**:

```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-05-09T10:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

### User Logout

Logs out a user by invalidating their current tokens.

- **URL**: `/api/v1/user/auth/logout`
- **Method**: `POST`
- **Auth Required**: Yes (User)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

### Refresh Token

Generates a new access token using a valid refresh token. This endpoint is shared between both admin and user domains.

- **URL**: `/api/v1/user/auth/refresh`
- **Method**: `POST`
- **Auth Required**: No (but requires a valid refresh token)

**Request Body**:

```json
{
  "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-05-09T11:30:00Z",
    "user_id": 1,
    "user_type": "user"
  }
}
```

### Get User Profile

Retrieves the profile of the currently authenticated user.

- **URL**: `/api/v1/user/auth/profile`
- **Method**: `GET`
- **Auth Required**: Yes (User)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2025-05-01T10:00:00Z",
    "updated_at": "2025-05-01T10:00:00Z",
    "last_login": "2025-05-09T10:00:00Z"
  }
}
```

### Change Password

Changes the password for the currently authenticated user.

- **URL**: `/api/v1/user/auth/change-password`
- **Method**: `POST`
- **Auth Required**: Yes (User)

**Request Body**:

```json
{
  "current_password": "SecurePass123!",
  "new_password": "EvenMoreSecure456!",
  "confirm_password": "EvenMoreSecure456!"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "Password updated successfully"
  }
}
```

## User Dashboard Endpoints

### Get User Dashboard

Retrieves the dashboard data for the currently authenticated user.

- **URL**: `/api/v1/user/dashboard`
- **Method**: `GET`
- **Auth Required**: Yes (User)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "Welcome to User Dashboard",
    "id": 1
  }
}
```

### Update User Profile

Updates the profile of the currently authenticated user.

- **URL**: `/api/v1/user/profile`
- **Method**: `PUT`
- **Auth Required**: Yes (User)

**Request Body**:

```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "created_at": "2025-05-01T10:00:00Z",
    "updated_at": "2025-05-09T11:00:00Z",
    "last_login": "2025-05-09T10:00:00Z"
  }
}
```

## Admin User Management Endpoints

These endpoints are for administrators to manage users.

### List Users

Retrieves a paginated list of users.

- **URL**: `/api/v1/admin/users`
- **Method**: `GET`
- **Auth Required**: Yes (Admin)

**Query Parameters**:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)
- `search`: Search term in name or email (optional)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "admin_id": 1,
      "created_at": "2025-05-01T10:00:00Z",
      "updated_at": "2025-05-01T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "admin_id": 1,
      "created_at": "2025-05-01T11:00:00Z",
      "updated_at": "2025-05-01T11:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "totalPages": 1
  }
}
```

### Create User

Creates a new user by an administrator.

- **URL**: `/api/v1/admin/users`
- **Method**: `POST`
- **Auth Required**: Yes (Admin)

**Request Body**:

```json
{
  "name": "New User",
  "email": "new.user@example.com"
}
```

**Response (201 Created)**:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": 3,
      "name": "New User",
      "email": "new.user@example.com",
      "admin_id": 1,
      "created_at": "2025-05-09T12:00:00Z",
      "updated_at": "2025-05-09T12:00:00Z"
    },
    "temporary_password": "Rand0mP@ssw0rd",
    "message": "User created successfully. Please inform the user to change their password after first login."
  }
}
```

### Get User

Retrieves a specific user by ID.

- **URL**: `/api/v1/admin/users/:id`
- **Method**: `GET`
- **Auth Required**: Yes (Admin)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com",
    "admin_id": 1,
    "created_at": "2025-05-01T10:00:00Z",
    "updated_at": "2025-05-01T10:00:00Z"
  }
}
```

### Update User

Updates a specific user.

- **URL**: `/api/v1/admin/users/:id`
- **Method**: `PUT`
- **Auth Required**: Yes (Admin)

**Request Body**:

```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "admin_id": 1,
    "created_at": "2025-05-01T10:00:00Z",
    "updated_at": "2025-05-09T13:00:00Z"
  }
}
```

### Delete User

Deletes a specific user.

- **URL**: `/api/v1/admin/users/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes (Admin)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "User deleted successfully"
  }
}
```

### Reset User Password

Resets a user's password. Admin can only reset passwords for users they created.

- **URL**: `/api/v1/admin/users/:id/reset-password`
- **Method**: `POST`
- **Auth Required**: Yes (Admin)

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "User password reset successfully",
    "new_password": "Temp@P4ssw0rd",
    "note": "Please provide this temporary password to the user and advise them to change it immediately after login."
  }
}
```

## Error Responses

### Authentication Error (401 Unauthorized)

```json
{
  "success": false,
  "error": "Unauthorized: user authentication required"
}
```

### Permission Error (403 Forbidden)

```json
{
  "success": false,
  "error": "Forbidden: you don't have permission to access this resource"
}
```

### Not Found Error (404 Not Found)

```json
{
  "success": false,
  "error": "User not found"
}
```

### Validation Error (400 Bad Request)

```json
{
  "success": false,
  "error": "Invalid input: email is required"
}
```

### Rate Limit Error (429 Too Many Requests)

```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later"
}
```

## Authentication Headers

All authenticated requests must include the JWT token in the Authorization header:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## Notes on User Types

1. **Self-Registered Users**: These users sign up through the registration endpoint and have no `admin_id` (or it is set to 0).

2. **Admin-Created Users**: These users are created by administrators and have an `admin_id` that refers to the admin who created them. Admin-created users can only be managed by the admin who created them.

3. **Password Handling**: 
   - Self-registered users set their own passwords during registration
   - Admin-created users receive a temporary password that they must change after first login