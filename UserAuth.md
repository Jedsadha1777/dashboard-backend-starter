# User Authentication System

This document explains the user authentication system that was added to the Dashboard Backend Starter application.

## Overview

The system now supports two separate authentication domains:

1. **Admin Authentication**: For administrative users who manage the dashboard
2. **User Authentication**: For regular users of the application

Each domain has its own registration, login, and management processes while sharing the same token infrastructure.

## User Authentication Endpoints

### Registration and Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | /api/v1/user/auth/register | Register a new user account |
| POST   | /api/v1/user/auth/login | Log in as a regular user |
| POST   | /api/v1/user/auth/logout | Log out (requires authentication) |
| POST   | /api/v1/user/auth/refresh | Refresh the access token |
| GET    | /api/v1/user/auth/profile | Get the current user profile |
| POST   | /api/v1/user/auth/change-password | Change user password |

### User Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /api/v1/user/dashboard | Get user dashboard data |

## Registration Process

1. User submits registration information (name, email, password, and password confirmation)
2. System validates the input and checks for duplicate emails
3. Password is hashed using bcrypt with a secure cost factor
4. User account is created with an initial token version
5. JWT access token and refresh token are generated
6. User receives tokens and can proceed to use the application

## Authentication Flow

1. **Login**: User provides email and password
2. Upon successful authentication:
   - JWT access token (30 minutes validity)
   - Refresh token (1 year validity) is issued
   - Token version is incremented to invalidate previous tokens
   - Last login timestamp is updated

3. **Token Usage**:
   - Access token is used for all authenticated requests
   - Include in the Authorization header: `Authorization: Bearer {token}`

4. **Token Refresh**:
   - When access token expires, use the refresh token to get a new access token
   - POST to `/api/v1/user/auth/refresh` with the refresh token
   - No need to login again

5. **Logout**:
   - Token version is incremented, invalidating all existing tokens
   - Requires authentication

## Security Features

1. **Password Security**:
   - Passwords are hashed using bcrypt
   - Minimum password length of 8 characters
   - Password confirmation required during registration

2. **Token Versioning**:
   - Each user has a token version number
   - Incremented on login, logout, and password change
   - Allows immediate invalidation of all previously issued tokens

3. **Refresh Token Security**:
   - Refresh tokens are stored in the database
   - Can be revoked individually or for all user sessions
   - Automatic cleanup of expired tokens
   - Protected against replay attacks

4. **Rate Limiting**:
   - Configurable rate limits for authentication endpoints
   - Prevents brute force attacks

## Example Registration Request

```bash
curl -X POST http://localhost:8080/api/v1/user/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "SecurePass123!",
    "confirm_password": "SecurePass123!"
  }'
```

Response:

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

## Example Login Request

```bash
curl -X POST http://localhost:8080/api/v1/user/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

Response:

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

## Example Change Password Request

```bash
curl -X POST http://localhost:8080/api/v1/user/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "current_password": "SecurePass123!",
    "new_password": "EvenMoreSecure456!",
    "confirm_password": "EvenMoreSecure456!"
  }'
```

Response:

```json
{
  "success": true,
  "data": {
    "message": "Password updated successfully"
  }
}
```

## Domain Separation

The system maintains a clear separation between admin and user domains:

1. **Different Authentication Endpoints**:
   - Admin: `/api/v1/auth/login`
   - User: `/api/v1/user/auth/login`

2. **Different Middleware Requirements**:
   - `AdminRequired()` for admin-only endpoints
   - `UserRequired()` for user-only endpoints
   - `SelfOrAdminRequired()` for endpoints that can be accessed by the owner or an admin

3. **Token Types**:
   - User tokens have `user_type: "user"`
   - Admin tokens have `user_type: "admin"`
   - The middleware checks the user type to enforce proper access control

4. **Database Storage**:
   - Admin users are stored in the `admin` table
   - Regular users are stored in the `user` table

## Integration with Frontend

When integrating with a frontend application:

1. Store tokens securely:
   - Access token: Store in memory or short-lived cookie
   - Refresh token: Store in an HTTP-only, secure cookie

2. Token refresh strategy:
   - Implement automatic refresh when access token expires
   - Use refresh token to obtain a new access token
   - Redirect to login if refresh fails

3. Protect against XSS and CSRF:
   - Never store tokens in localStorage
   - Use HTTP-only cookies for refresh tokens
   - Implement proper CSRF protection

## Error Handling

The authentication system returns standardized error responses:

| Status Code | Type of Error |
|-------------|---------------|
| 400 | Invalid input, validation errors |
| 401 | Authentication failed, invalid credentials |
| 403 | Authorization failed, insufficient permissions |
| 429 | Too many requests (rate limiting) |
| 500 | Server error |

Example error response:

```json
{
  "success": false,
  "error": "Invalid email or password"
}
```