# LUMEN API Documentation

## Overview

LUMEN API is a production-ready Go backend service with Supabase PostgreSQL integration for habit tracking, task management, and daily logging.

**Base URL**: `http://localhost:8080`

**Version**: 1.0.0

## Authentication

All API endpoints (except health checks) require JWT authentication.

Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

### Success Response
```json
{
  "data": {},
  "count": 0
}
```

### Error Response
```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message"
}
```

## Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Resource deleted successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation failed
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

## Endpoints

### Health Check

#### GET /health

Check API health status.

**Response**
```json
{
  "status": "ok",
  "timestamp": "2025-11-13T10:00:00Z",
  "service": "lumen-api",
  "version": "1.0.0",
  "database": "healthy"
}
```

#### GET /ready

Check if service is ready to accept requests.

**Response**
```json
{
  "status": "ready"
}
```

#### GET /metrics

Get database connection metrics.

**Response**
```json
{
  "database": {
    "acquired_conns": 2,
    "idle_conns": 3,
    "total_conns": 5,
    "max_conns": 25
  },
  "timestamp": "2025-11-13T10:00:00Z"
}
```

---

### Habits

#### GET /api/habits

Get all habits for authenticated user.

**Headers**
```
Authorization: Bearer <token>
```

**Response**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "name": "Morning Exercise",
      "color": "#FF5733",
      "icon": "dumbbell",
      "frequency": "daily",
      "target_count": 1,
      "is_active": true,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ],
  "count": 1
}
```

#### POST /api/habits

Create a new habit.

**Request Body**
```json
{
  "name": "Morning Exercise",
  "color": "#FF5733",
  "icon": "dumbbell",
  "frequency": "daily",
  "target_count": 1
}
```

**Validation Rules**
- `name`: required, 1-100 characters
- `color`: required, valid hex color
- `icon`: required, 1-50 characters
- `frequency`: required, one of: `daily`, `weekly`, `monthly`
- `target_count`: required, 1-100

**Response** (201 Created)
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "Morning Exercise",
  "color": "#FF5733",
  "icon": "dumbbell",
  "frequency": "daily",
  "target_count": 1,
  "is_active": true,
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### GET /api/habits/:id

Get a specific habit by ID.

**Parameters**
- `id` (path): Habit UUID

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "Morning Exercise",
  "color": "#FF5733",
  "icon": "dumbbell",
  "frequency": "daily",
  "target_count": 1,
  "is_active": true,
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### PUT /api/habits/:id

Update a habit.

**Parameters**
- `id` (path): Habit UUID

**Request Body** (all fields optional)
```json
{
  "name": "Evening Exercise",
  "color": "#33FF57",
  "icon": "running",
  "frequency": "weekly",
  "target_count": 3,
  "is_active": false
}
```

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "Evening Exercise",
  "color": "#33FF57",
  "icon": "running",
  "frequency": "weekly",
  "target_count": 3,
  "is_active": false,
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:30:00Z"
}
```

#### DELETE /api/habits/:id

Delete a habit.

**Parameters**
- `id` (path): Habit UUID

**Response** (204 No Content)

---

### Tasks

#### GET /api/tasks

Get all tasks for authenticated user with optional filters.

**Query Parameters**
- `horizon` (optional): Filter by horizon - `now`, `next`, `later`, `someday`
- `status` (optional): Filter by status - `todo`, `in_progress`, `done`, `archived`
- `priority` (optional): Filter by priority - `low`, `medium`, `high`, `urgent`

**Example**: `/api/tasks?horizon=now&status=todo`

**Response**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "title": "Complete project proposal",
      "description": "Write detailed project proposal for Q4",
      "horizon": "now",
      "priority": "high",
      "status": "in_progress",
      "due_date": "2025-11-15T00:00:00Z",
      "completed_at": null,
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ],
  "count": 1,
  "filter": {
    "horizon": "now",
    "status": "todo"
  }
}
```

#### POST /api/tasks

Create a new task.

**Request Body**
```json
{
  "title": "Complete project proposal",
  "description": "Write detailed project proposal for Q4",
  "horizon": "now",
  "priority": "high",
  "due_date": "2025-11-15T00:00:00Z"
}
```

**Validation Rules**
- `title`: required, 1-200 characters
- `description`: optional, max 1000 characters
- `horizon`: required, one of: `now`, `next`, `later`, `someday`
- `priority`: required, one of: `low`, `medium`, `high`, `urgent`
- `due_date`: optional, ISO 8601 datetime

**Response** (201 Created)
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Complete project proposal",
  "description": "Write detailed project proposal for Q4",
  "horizon": "now",
  "priority": "high",
  "status": "todo",
  "due_date": "2025-11-15T00:00:00Z",
  "completed_at": null,
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### GET /api/tasks/:id

Get a specific task by ID.

**Parameters**
- `id` (path): Task UUID

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Complete project proposal",
  "description": "Write detailed project proposal for Q4",
  "horizon": "now",
  "priority": "high",
  "status": "in_progress",
  "due_date": "2025-11-15T00:00:00Z",
  "completed_at": null,
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### PUT /api/tasks/:id

Update a task.

**Parameters**
- `id` (path): Task UUID

**Request Body** (all fields optional)
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "horizon": "next",
  "priority": "medium",
  "status": "done",
  "due_date": "2025-11-20T00:00:00Z"
}
```

**Note**: When status is changed to `done`, `completed_at` is automatically set.

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "title": "Updated title",
  "description": "Updated description",
  "horizon": "next",
  "priority": "medium",
  "status": "done",
  "due_date": "2025-11-20T00:00:00Z",
  "completed_at": "2025-11-13T10:30:00Z",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:30:00Z"
}
```

#### DELETE /api/tasks/:id

Delete a task.

**Parameters**
- `id` (path): Task UUID

**Response** (204 No Content)

---

### Daily Logs

#### POST /api/daily-log

Create or update a daily log for a specific date.

**Request Body**
```json
{
  "date": "2025-11-13T00:00:00Z",
  "morning_routine": true,
  "evening_routine": false,
  "water_intake": 8,
  "sleep_hours": 7.5,
  "energy_level": 4,
  "mood_rating": 5,
  "productivity_rating": 4,
  "notes": "Great day, very productive"
}
```

**Validation Rules**
- `date`: required, ISO 8601 date
- `morning_routine`: boolean
- `evening_routine`: boolean
- `water_intake`: 0-20 glasses
- `sleep_hours`: 0-24 hours
- `energy_level`: 1-5 rating
- `mood_rating`: 1-5 rating
- `productivity_rating`: 1-5 rating
- `notes`: optional, max 1000 characters

**Response** (201 Created)
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "date": "2025-11-13T00:00:00Z",
  "morning_routine": true,
  "evening_routine": false,
  "water_intake": 8,
  "sleep_hours": 7.5,
  "energy_level": 4,
  "mood_rating": 5,
  "productivity_rating": 4,
  "notes": "Great day, very productive",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### GET /api/daily-log/:date

Get daily log for a specific date.

**Parameters**
- `date` (path): Date in format YYYY-MM-DD

**Example**: `/api/daily-log/2025-11-13`

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "date": "2025-11-13T00:00:00Z",
  "morning_routine": true,
  "evening_routine": false,
  "water_intake": 8,
  "sleep_hours": 7.5,
  "energy_level": 4,
  "mood_rating": 5,
  "productivity_rating": 4,
  "notes": "Great day, very productive",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:00:00Z"
}
```

#### GET /api/daily-log

Get daily logs for a date range.

**Query Parameters**
- `start_date` (required): Start date in format YYYY-MM-DD
- `end_date` (required): End date in format YYYY-MM-DD

**Example**: `/api/daily-log?start_date=2025-11-01&end_date=2025-11-13`

**Response**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "date": "2025-11-13T00:00:00Z",
      "morning_routine": true,
      "evening_routine": false,
      "water_intake": 8,
      "sleep_hours": 7.5,
      "energy_level": 4,
      "mood_rating": 5,
      "productivity_rating": 4,
      "notes": "Great day, very productive",
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  ],
  "count": 1,
  "start_date": "2025-11-01",
  "end_date": "2025-11-13"
}
```

#### PUT /api/daily-log/:date

Update a daily log for a specific date.

**Parameters**
- `date` (path): Date in format YYYY-MM-DD

**Request Body** (all fields optional)
```json
{
  "morning_routine": true,
  "evening_routine": true,
  "water_intake": 10,
  "sleep_hours": 8.0,
  "energy_level": 5,
  "mood_rating": 5,
  "productivity_rating": 5,
  "notes": "Updated notes"
}
```

**Response**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "date": "2025-11-13T00:00:00Z",
  "morning_routine": true,
  "evening_routine": true,
  "water_intake": 10,
  "sleep_hours": 8.0,
  "energy_level": 5,
  "mood_rating": 5,
  "productivity_rating": 5,
  "notes": "Updated notes",
  "created_at": "2025-11-13T10:00:00Z",
  "updated_at": "2025-11-13T10:30:00Z"
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Limit**: 100 requests per minute per IP/user
- **Response**: 429 Too Many Requests when limit exceeded

## Error Codes

| Code | Description |
|------|-------------|
| `BAD_REQUEST` | Invalid request parameters |
| `UNAUTHORIZED` | Missing or invalid authentication |
| `FORBIDDEN` | Insufficient permissions |
| `NOT_FOUND` | Resource not found |
| `CONFLICT` | Resource conflict (duplicate) |
| `VALIDATION_ERROR` | Request validation failed |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `DATABASE_ERROR` | Database operation failed |
| `INTERNAL_SERVER_ERROR` | Unexpected server error |

## Architecture

### Repository Pattern
- Clean separation of concerns
- Database abstraction layer
- Interface-based design for testability

### Middleware Stack
1. Recovery - Panic recovery
2. RequestLogger - Request/response logging
3. CORS - Cross-origin resource sharing
4. RateLimit - Rate limiting protection
5. Authentication - JWT validation

### Database Connection
- Connection pooling (5-25 connections)
- Health checks every minute
- Automatic reconnection
- Connection lifecycle management

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL (via Supabase)
- Make (optional)

### Setup

1. Copy environment variables:
```bash
cp .env.example .env
```

2. Update `.env` with your Supabase credentials

3. Install dependencies:
```bash
go mod download
```

4. Run the server:
```bash
go run cmd/server/main.go
```

Or with Make:
```bash
make dev
```

### Build

```bash
make build
./bin/server
```

### Docker

```bash
make docker-build
make docker-run
```

## Production Deployment

### Environment Variables

Required variables for production:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `SUPABASE_URL`, `SUPABASE_KEY`
- `JWT_SECRET` (minimum 32 characters)
- `PORT` (default: 8080)
- `GIN_MODE=release`

### Health Checks

Use `/health` endpoint for load balancer health checks:
- Returns 200 when healthy
- Returns 503 when database is unavailable

### Logging

Structured JSON logging with:
- Request IDs for tracing
- Contextual fields
- Error stack traces
- Performance metrics

### Security

- JWT authentication on all API endpoints
- Rate limiting (100 req/min)
- CORS configuration
- SQL injection prevention (parameterized queries)
- Input validation on all endpoints
- Secure headers
