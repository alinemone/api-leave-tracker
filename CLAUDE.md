# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Leave Tracker application built with the Goravel framework (Go-based Laravel-style framework). It manages employee leave requests with features for creating, listing, and reporting on leave records. The application integrates with GAP (a messaging service) for notifications and uses PostgreSQL for data persistence.

## Development Commands

### Running the Application
```bash
# Run with hot reload (using Air)
go run . artisan artisan:air

# Generate application key
go run . artisan key:generate

# Run without hot reload
go run main.go
```

### Database
```bash
# Run migrations
go run . artisan migrate

# Create a new migration
go run . artisan make:migration <migration_name>
```

### Testing
```bash
# Run all tests
go test ./tests/...

# Run feature tests
go test ./tests/feature/...

# Run with verbose output
go test -v ./tests/...
```

### Building
```bash
# Build the application
go build -o ./storage/temp/main.exe .
```

## Architecture

### Core Framework
- **Goravel Framework**: Laravel-style MVC framework for Go
- **Database**: PostgreSQL with GORM ORM
- **Router**: Gin-based HTTP routing via goravel/gin
- **Authentication**: JWT-based authentication

### Project Structure

**Routes** (`routes/`):
- `api.go`: API endpoints at `/api/v1` prefix
- `web.go`: Web routes
- `grpc.go`: gRPC service definitions

**Application Layer** (`app/`):
- `models/`: Database models (User, Leave)
- `repositories/`: Data access layer implementing interfaces
- `interfaces/`: Repository interfaces (UserInterface, LeaveInterface)
- `http/controllers/`: HTTP request handlers
- `http/requests/`: Form request validation
- `http/responses/`: Response formatters
- `http/middleware/`: HTTP middleware (auth)
- `events/`: Event definitions
- `listeners/`: Event listeners for notifications
- `rules/`: Custom validation rules
- `providers/`: Service providers for dependency injection
- `helpers/`: Utility functions

**Services** (`services/`):
- `gap/`: GAP messaging SDK client for notifications

**Configuration** (`config/`):
- Database, app, cache, and other configuration files

**Database** (`database/`):
- `migrations/`: Database migration files

### Key Architectural Patterns

**Repository Pattern**: Data access is abstracted through interfaces defined in `app/interfaces/` and implemented in `app/repositories/`. Controllers depend on interfaces, not concrete implementations.

**Service Provider Pattern**: Dependencies are registered in service providers (`app/providers/`) and resolved via Goravel's IoC container.

**Event-Driven Notifications**: Leave creation/deletion triggers events (`app/events/leave/`) that are handled by listeners (`app/listeners/leave/`) to send GAP notifications.

**Request Validation**: Form requests (`app/http/requests/`) handle authorization, validation rules, and data preparation before reaching controllers.

### Leave Management System

**Leave Types**:
- `daily`: Full-day leave
- `hourly`: Partial-day leave with specific time ranges

**Validation Rules**:
- `leave_overlap`: Custom rule in `app/rules/leave/leave_overlap.go` prevents overlapping leaves
  - Daily leaves: No other leave on same day
  - Hourly leaves: No time overlap on same day
- `same_day`: Custom rule ensures start and end are on same day for hourly leaves

**Leave Model** (`app/models/leave.go`):
- Uses soft deletes
- Fields: UserId, Type, StartAt, EndAt

### Authentication Flow
- JWT tokens used for API authentication
- Auth middleware in `app/http/middleware/auth.go`
- CurrentUser helper in `app/helpers/helpers.go` extracts authenticated user from context

### GAP Integration
The `services/gap/gap_sdk.go` client handles:
- Sending messages to rooms
- Sending direct messages to users
- User info lookup

Configuration in `.env`:
- `GAP_BASE_URL`: GAP service URL
- `GAP_AUTH_TOKEN`: Authentication token
- `GAP_USER_ID`: Bot user ID
- `GAP_ROOM_ID`: Default room for notifications

## Environment Configuration

Copy `.env.example` to `.env` and configure:
- Database connection (DB_HOST, DB_PORT, DB_DATABASE, DB_USERNAME, DB_PASSWORD)
- JWT secret (JWT_SECRET)
- GAP credentials for notifications
- Sentry DSN for error tracking

## Dependencies

Main dependencies in `go.mod`:
- `github.com/goravel/framework` - Core framework
- `github.com/goravel/gin` - HTTP routing
- `github.com/goravel/postgres` - PostgreSQL driver
- `github.com/getsentry/sentry-go` - Error tracking
- `gorm.io/gorm` - ORM
- `golang.org/x/crypto` - Password hashing

## Common Workflows

### Adding a New Entity
1. Create model in `app/models/`
2. Define interface in `app/interfaces/`
3. Implement repository in `app/repositories/`
4. Register in `app/providers/app_service_provider.go`
5. Create controller in `app/http/controllers/`
6. Add routes in `routes/api.go`
7. Create request validators in `app/http/requests/`
8. Create response formatters in `app/http/responses/`

### Running Tests for a Single File
```bash
go test ./tests/feature/<test_file>.go
```

## Notes

- The application uses Persian text in validation messages and reports
- Air configuration in `.air.toml` enables hot reload during development
- Sentry is initialized in `main.go` for error tracking
- Recent changes focus on leave validation and overlap checking (see git log)
