# Article Management API

This document outlines the API endpoints for managing articles in the Admin Dashboard.

## Authentication

All article management endpoints require admin authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer YOUR_JWT_TOKEN
```

## Endpoints

### Create Article

Creates a new article.

- **URL**: `/api/v1/admin/articles`
- **Method**: `POST`
- **Auth Required**: Yes

**Request Body**:

```json
{
  "title": "Article Title",
  "content": "Article content goes here...",
  "slug": "article-slug",
  "summary": "A brief summary of the article",
  "status": "draft",
  "published_at": "2025-05-07T10:00:00Z"
}
```

Notes:
- If `slug` is not provided, one will be generated from the title
- Valid status values: `draft`, `published`, `archived` (defaults to `draft`)
- `published_at` is optional (ISO 8601 format)

**Response (201 Created)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Article Title",
    "content": "Article content goes here...",
    "slug": "article-slug",
    "summary": "A brief summary of the article",
    "status": "draft",
    "published_at": null,
    "admin_id": 1,
    "created_at": "2025-05-07T10:00:00Z",
    "updated_at": "2025-05-07T10:00:00Z"
  }
}
```

### List Articles

Retrieves a paginated list of articles with optional filtering.

- **URL**: `/api/v1/admin/articles`
- **Method**: `GET`
- **Auth Required**: Yes

**Query Parameters**:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `search`: Search term (optional)
- `status`: Filter by status (optional) - `draft`, `published`, or `archived`

**Response (200 OK)**:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Article Title",
      "content": "Article content goes here...",
      "slug": "article-slug",
      "summary": "A brief summary of the article",
      "status": "draft",
      "published_at": null,
      "admin_id": 1,
      "admin": {
        "id": 1,
        "email": "admin@example.com"
      },
      "created_at": "2025-05-07T10:00:00Z",
      "updated_at": "2025-05-07T10:00:00Z"
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

### Get Article

Retrieves a specific article by ID.

- **URL**: `/api/v1/admin/articles/:id`
- **Method**: `GET`
- **Auth Required**: Yes

**URL Parameters**:

- `id`: Article ID

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Article Title",
    "content": "Article content goes here...",
    "slug": "article-slug",
    "summary": "A brief summary of the article",
    "status": "draft",
    "published_at": null,
    "admin_id": 1,
    "admin": {
      "id": 1,
      "email": "admin@example.com"
    },
    "created_at": "2025-05-07T10:00:00Z",
    "updated_at": "2025-05-07T10:00:00Z"
  }
}
```

### Update Article

Updates an existing article.

- **URL**: `/api/v1/admin/articles/:id`
- **Method**: `PUT`
- **Auth Required**: Yes

**URL Parameters**:

- `id`: Article ID

**Request Body**:

```json
{
  "title": "Updated Article Title",
  "content": "Updated content goes here...",
  "slug": "updated-article-slug",
  "summary": "Updated summary of the article",
  "status": "published",
  "published_at": "2025-05-08T15:30:00Z"
}
```

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Updated Article Title",
    "content": "Updated content goes here...",
    "slug": "updated-article-slug",
    "summary": "Updated summary of the article",
    "status": "published",
    "published_at": "2025-05-08T15:30:00Z",
    "admin_id": 1,
    "admin": {
      "id": 1,
      "email": "admin@example.com"
    },
    "created_at": "2025-05-07T10:00:00Z",
    "updated_at": "2025-05-08T15:30:00Z"
  }
}
```

Note: Only the admin who created the article can update it.

### Delete Article

Deletes an article.

- **URL**: `/api/v1/admin/articles/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes

**URL Parameters**:

- `id`: Article ID

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "message": "Article deleted successfully"
  }
}
```

Note: Only the admin who created the article can delete it.

### Publish Article

Sets an article's status to 'published' and sets the published_at timestamp.

- **URL**: `/api/v1/admin/articles/:id/publish`
- **Method**: `POST`
- **Auth Required**: Yes

**URL Parameters**:

- `id`: Article ID

**Response (200 OK)**:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Article Title",
    "content": "Article content goes here...",
    "slug": "article-slug",
    "summary": "A brief summary of the article",
    "status": "published",
    "published_at": "2025-05-08T15:45:00Z",
    "admin_id": 1,
    "admin": {
      "id": 1,
      "email": "admin@example.com"
    },
    "created_at": "2025-05-07T10:00:00Z",
    "updated_at": "2025-05-08T15:45:00Z"
  }
}
```

Note: Only the admin who created the article can publish it.

## Error Responses

### Authentication Error (401 Unauthorized)

```json
{
  "success": false,
  "error": "Unauthorized: admin authentication required"
}
```

### Permission Error (403 Forbidden)

```json
{
  "success": false,
  "error": "Failed to update article: you don't have permission to update this article"
}
```

### Not Found Error (404 Not Found)

```json
{
  "success": false,
  "error": "Article not found"
}
```

### Validation Error (400 Bad Request)

```json
{
  "success": false,
  "error": "Invalid input: title is required"
}
```

### Server Error (500 Internal Server Error)

```json
{
  "success": false,
  "error": "Failed to create article: database connection error"
}
```