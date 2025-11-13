# LUMEN Backend - Go Service

Production-ready Go backend service for the LUMEN personal operating system.

## Architecture

### Project Structure

```
backend-go/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── api/            # HTTP handlers and routes
│   ├── domain/         # Business logic and domain models
│   ├── repository/     # Data access layer
│   └── middleware/     # HTTP middleware
├── pkg/                # Public packages (can be imported by external apps)
│   ├── config/         # Configuration management
│   └── logger/         # Logging utilities
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Design Principles

1. **Clean Architecture**: Separation of concerns with clear boundaries
2. **Dependency Injection**: Interface-based design for testability
3. **Repository Pattern**: Abstract data access layer
4. **Middleware Pipeline**: Composable request processing
5. **Structured Logging**: Production-ready observability

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+ (via Supabase)
- Redis (optional, for caching)
- Make (for build automation)

### Installation

1. Install dependencies:
```bash
make install
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your actual values
```

3. Build the application:
```bash
make build
```

### Development

Run the development server:
```bash
make run
```

Or with hot reload (requires [air](https://github.com/cosmtrek/air)):
```bash
make dev
```

### Testing

Run all tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

### Code Quality

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

Run go vet:
```bash
make vet
```

## API Documentation

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Authentication

All API routes (except `/health`) require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options.

### Key Configuration Areas

- **Application**: Port, environment, name
- **Database**: Connection string, pool settings
- **Redis**: Cache configuration
- **JWT**: Secret key, token expiry
- **CORS**: Allowed origins and methods
- **Logging**: Level, format, output
- **Rate Limiting**: Request limits and windows

## Deployment

### Docker

Build Docker image:
```bash
make docker-build
```

Run Docker container:
```bash
make docker-run
```

### Production Checklist

- [ ] Set `APP_ENV=production`
- [ ] Use strong `JWT_SECRET` (32+ characters)
- [ ] Configure `DATABASE_URL` with SSL
- [ ] Set appropriate `CORS_ALLOWED_ORIGINS`
- [ ] Enable rate limiting
- [ ] Set `LOG_LEVEL=info` or `warn`
- [ ] Configure monitoring and alerting
- [ ] Set up automated backups

## Best Practices

### Error Handling

Always return structured errors:

```go
type AppError struct {
    Code        string `json:"code"`
    Message     string `json:"message"`
    UserMessage string `json:"user_message"`
    StatusCode  int    `json:"status_code"`
}
```

### Logging

Use structured logging with context:

```go
logger.Info("habit created",
    zap.String("user_id", userID),
    zap.String("habit_id", habitID),
)
```

### Database Transactions

Always use context and handle errors:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()

// Perform operations...

return tx.Commit()
```

## Testing Strategy

### Unit Tests
- Test business logic in isolation
- Mock external dependencies
- Aim for >90% coverage

### Integration Tests
- Test with real database (test container)
- Verify API endpoints end-to-end
- Test middleware and auth flows

### Load Tests
- Use tools like `hey` or `vegeta`
- Test rate limiting
- Verify connection pooling

## Troubleshooting

### Common Issues

**Database connection fails:**
- Verify `DATABASE_URL` is correct
- Check network connectivity to Supabase
- Ensure SSL mode is properly configured

**JWT verification fails:**
- Verify `JWT_SECRET` matches frontend
- Check token expiry settings
- Ensure clock synchronization

**CORS errors:**
- Add frontend URL to `CORS_ALLOWED_ORIGINS`
- Verify allowed methods and headers
- Check preflight request handling

## Contributing

1. Follow Go best practices and conventions
2. Write tests for new features
3. Update documentation
4. Run linter before committing
5. Keep functions small and focused

## License

MIT
