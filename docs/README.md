# LUMEN Backend - Go API

Production-ready Go backend service for LUMEN personal operating system with Supabase PostgreSQL integration.

## Overview

LUMEN Backend is a RESTful API built with Go 1.21+, Gin framework, and Supabase PostgreSQL. It provides endpoints for habit tracking, task management, and daily logging with robust authentication, error handling, and logging.

## Features

- RESTful API with Gin framework
- PostgreSQL database via Supabase
- JWT authentication
- Comprehensive error handling
- Structured logging with Zap
- CORS middleware
- Rate limiting (100 req/min)
- Health check endpoints
- Database connection pooling
- Repository pattern architecture
- Docker support
- Production-ready configuration

## Tech Stack

- **Go**: 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (Supabase)
- **Driver**: pgx v5
- **Logging**: Zap
- **Environment**: godotenv

## Project Structure

```
backend-go/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── handlers/             # HTTP request handlers
│   │   ├── health_handler.go
│   │   ├── habit_handler.go
│   │   ├── task_handler.go
│   │   └── daily_log_handler.go
│   ├── middleware/           # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── rate_limit.go
│   ├── models/               # Data models
│   │   ├── habit.go
│   │   ├── task.go
│   │   ├── daily_log.go
│   │   ├── user.go
│   │   └── errors.go
│   └── repository/           # Database layer
│       ├── database.go
│       ├── habit_repository.go
│       ├── task_repository.go
│       └── daily_log_repository.go
├── pkg/
│   ├── errors/               # Custom error types
│   │   └── errors.go
│   └── logger/               # Logging utilities
│       └── logger.go
├── docs/
│   ├── API.md                # API documentation
│   ├── DEPLOYMENT.md         # Deployment guide
│   └── README.md             # This file
├── .env.example              # Example environment variables
├── .gitignore
├── Dockerfile
├── Makefile
└── go.mod
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Supabase account)
- Make (optional)

### Installation

1. Clone the repository:
```bash
cd backend-go
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Update `.env` with your Supabase credentials:
```env
PORT=8080
GIN_MODE=debug
DB_HOST=db.your-project.supabase.co
DB_PASSWORD=your-db-password
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key
JWT_SECRET=your-secure-jwt-secret
```

4. Install dependencies:
```bash
go mod download
```

5. Run the server:
```bash
make dev
# or
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Verify Installation

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "ok",
#   "timestamp": "2025-11-13T10:00:00Z",
#   "service": "lumen-api",
#   "version": "1.0.0",
#   "database": "healthy"
# }
```

## API Endpoints

### Health Checks

- `GET /health` - Health status
- `GET /ready` - Readiness check
- `GET /metrics` - Database metrics

### Habits

- `GET /api/habits` - Get all habits
- `POST /api/habits` - Create habit
- `GET /api/habits/:id` - Get habit by ID
- `PUT /api/habits/:id` - Update habit
- `DELETE /api/habits/:id` - Delete habit

### Tasks

- `GET /api/tasks` - Get all tasks (with filters)
- `POST /api/tasks` - Create task
- `GET /api/tasks/:id` - Get task by ID
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task

### Daily Logs

- `POST /api/daily-log` - Create/update daily log
- `GET /api/daily-log/:date` - Get log by date
- `GET /api/daily-log` - Get logs by date range
- `PUT /api/daily-log/:date` - Update daily log

See [API.md](./API.md) for complete API documentation.

## Development

### Available Commands

```bash
make help          # Show all available commands
make install       # Install dependencies
make build         # Build the application
make run           # Build and run
make dev           # Run in development mode
make test          # Run tests with coverage
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make lint          # Run linter
make format        # Format code
```

### Running Tests

```bash
make test
```

This will:
- Run all tests with race detection
- Generate coverage report
- Create HTML coverage report

### Code Style

- Follow Go standard naming conventions
- Use interfaces for dependency injection
- Keep functions under 100 lines
- Handle all errors explicitly
- Use structured logging

### Adding a New Endpoint

1. Create model in `internal/models/`
2. Create repository interface and implementation in `internal/repository/`
3. Create handler in `internal/handlers/`
4. Register route in `cmd/server/main.go`
5. Add tests
6. Update API documentation

Example:

```go
// 1. Model (internal/models/example.go)
type Example struct {
    ID        uuid.UUID `json:"id" db:"id"`
    UserID    uuid.UUID `json:"user_id" db:"user_id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// 2. Repository (internal/repository/example_repository.go)
type ExampleRepository interface {
    Create(ctx context.Context, example *Example) error
    GetByID(ctx context.Context, id, userID uuid.UUID) (*Example, error)
}

// 3. Handler (internal/handlers/example_handler.go)
type ExampleHandler struct {
    repo repository.ExampleRepository
}

func (h *ExampleHandler) Create(c *gin.Context) {
    // Implementation
}

// 4. Register route (cmd/server/main.go)
exampleHandler := handlers.NewExampleHandler(exampleRepo)
api.POST("/examples", exampleHandler.Create)
```

## Configuration

### Environment Variables

See `.env.example` for all available configuration options.

Required variables:
- `DB_HOST`, `DB_PASSWORD` - Database connection
- `SUPABASE_URL`, `SUPABASE_KEY` - Supabase configuration
- `JWT_SECRET` - JWT signing secret (min 32 chars)

Optional variables:
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (debug/release)
- `LOG_LEVEL` - Log level (debug/info/warn/error)
- `ALLOWED_ORIGINS` - CORS allowed origins

### Database Configuration

Connection pool settings (in `repository/database.go`):
```go
config.MaxConns = 25              // Maximum connections
config.MinConns = 5               // Minimum idle connections
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
config.HealthCheckPeriod = time.Minute
```

## Deployment

### Docker

Build and run with Docker:

```bash
make docker-build
make docker-run
```

Or manually:

```bash
docker build -t lumen-api:latest .
docker run -p 8080:8080 --env-file .env lumen-api:latest
```

### Production

For production deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md).

Supported platforms:
- Railway
- Render
- Google Cloud Run
- AWS ECS
- Any Docker-compatible platform

## Architecture

### Repository Pattern

Clean separation of concerns:
- **Handlers**: HTTP request/response handling
- **Repository**: Database operations
- **Models**: Data structures and validation
- **Middleware**: Cross-cutting concerns

### Dependency Injection

All dependencies injected through interfaces:

```go
type HabitRepository interface {
    Create(ctx context.Context, habit *Habit) error
    GetByID(ctx context.Context, id, userID uuid.UUID) (*Habit, error)
}

// Implementation
type habitRepository struct {
    db *Database
}

// Injection
habitRepo := repository.NewHabitRepository(db)
habitHandler := handlers.NewHabitHandler(habitRepo)
```

### Error Handling

Consistent error responses:

```go
type AppError struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    StatusCode int    `json:"-"`
    Err        error  `json:"-"`
}
```

All endpoints return either data or AppError.

### Middleware Stack

Request flow:
1. Recovery (panic handling)
2. RequestLogger (request/response logging)
3. CORS (cross-origin resource sharing)
4. RateLimit (rate limiting)
5. Authentication (JWT validation)
6. Handler

## Security

- JWT authentication on all API endpoints
- Row Level Security (RLS) in Supabase
- Rate limiting (100 requests/minute)
- CORS configuration
- SQL injection prevention (parameterized queries)
- Input validation
- Structured logging (no sensitive data)

## Performance

- Connection pooling (5-25 connections)
- Efficient database queries with indexes
- Minimal memory allocation
- Fast JSON serialization
- Health check caching

## Monitoring

### Logs

Structured JSON logging with:
- Request IDs for tracing
- Timestamp
- Log level
- Contextual fields
- Error stack traces

Example log:
```json
{
  "level": "info",
  "timestamp": "2025-11-13T10:00:00Z",
  "request_id": "uuid",
  "method": "GET",
  "path": "/api/habits",
  "status": 200,
  "latency": "15ms",
  "message": "Request completed"
}
```

### Metrics

Database connection pool metrics available at `/metrics`:
```json
{
  "database": {
    "acquired_conns": 2,
    "idle_conns": 3,
    "total_conns": 5,
    "max_conns": 25
  }
}
```

## Troubleshooting

### Common Issues

**Database connection failed:**
- Verify Supabase credentials
- Check network connectivity
- Ensure SSL mode is correct

**Port already in use:**
```bash
lsof -i :8080
kill -9 <PID>
```

**JWT authentication failing:**
- Verify JWT_SECRET is set
- Check token format
- Ensure token hasn't expired

See [DEPLOYMENT.md](./DEPLOYMENT.md) for more troubleshooting tips.

## Contributing

1. Follow Go code conventions
2. Write tests for new features
3. Update documentation
4. Keep functions small (<100 lines)
5. Handle all errors
6. Use structured logging

## License

This project is part of LUMEN personal operating system.

## Support

- API Documentation: [API.md](./API.md)
- Deployment Guide: [DEPLOYMENT.md](./DEPLOYMENT.md)
- Issues: Create an issue in the repository
