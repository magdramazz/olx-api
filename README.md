# OLX API

A RESTful API service for managing classified listings similar to OLX, built with Go (Golang).

## Overview

This project provides a simple but robust API for creating, listing, and deleting classified advertisements. It features a PostgreSQL database backend, proper error handling, request tracing, and migration support.

## Features

- **RESTful Endpoints**: Create, list, and delete listings
- **Database Persistence**: PostgreSQL with connection pooling
- **Database Migrations**: Version-controlled schema changes
- **Request Tracing**: Unique request IDs for debugging
- **Input Validation**: Server-side validation of listing data
- **Proper Error Handling**: Consistent JSON error responses
- **Health Check**: Endpoint for monitoring service status
- **Environment Configuration**: Configurable via environment variables
- **Logging**: Structured logging with slog

## API Endpoints

All endpoints return JSON responses.

### Health Check
```
GET /healthz
```
Returns service status.

### Listings
```
GET /listings
```
Retrieve a list of recent listings (limited to 100 most recent).

```
POST /listings
```
Create a new listing.

**Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "price": "string (non-negative integer)",
  "city": "string"
}
```

**Response:**
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "price": "string",
  "city": "string",
  "created_at": "timestamp"
}
```

```
DELETE /listings/{id}
```
Delete a listing by its UUID.

## Setup and Installation

### Prerequisites
- Go 1.26+ (tested with 1.26.1)
- PostgreSQL database
- Git

### Installation Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd olx-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   Copy `.env.example` to `.env` and fill in the required values:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your configuration:
   ```env
   PORT=8080
   ENV=development
   DATABASE_URL="postgresql://username:password@host:port/database?sslmode=..."
   ```

4. **Run database migrations**
   ```bash
   go run ./cmd/migrate up
   ```

5. **Build and run the application**
   ```bash
   make build
   ./bin/api
   ```
   Or run directly:
   ```bash
   make run
   ```

## Database Migrations

Migrations are managed with [golang-migrate](https://github.com/golang-migrate/migrate). To manage migrations:

```bash
# Apply all pending migrations
go run ./cmd/migrate up

# Rollback the last migration
go run ./cmd/migrate down
```

Seed data can be found in the `seeds/` directory and can be loaded manually into your database.

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `PORT` | Port for the HTTP server | Yes |
| `ENV` | Environment (development, production, etc.) | Yes |
| `DATABASE_URL` | PostgreSQL connection string | Yes |

## Project Structure

```
olx-api/
├── .github/               # GitHub workflows and configs
├── cmd/                   # Application entry points
│   ├── api/               # Main API server
│   └── migrate/           # Migration CLI tool
├── internal/              # Private application code
│   ├── config/            # Configuration loading
│   ├── db/                # Database connection
│   ├── handlers/          # HTTP request handlers
│   ├── httpx/             # HTTP utilities (error handling)
│   └── middleware/        # HTTP middleware
├── migrations/            # Database migration files
├── seeds/                 # Sample data for database
├── testsprite-tests/      # Test files (if any)
├── .env.example           # Example environment file
├── .gitignore             # Git ignore rules
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile               # Build and run commands
└── README.md              # This file
```

## Technologies Used

- **Language**: Go 1.26.1
- **Framework**: Standard library `net/http`
- **Database**: PostgreSQL with [pgx/v5](https://github.com/jackc/pgx/v5)
- **Migrations**: [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate)
- **Environment**: [godotenv](https://github.com/joho/godotenv)
- **UUIDs**: [google/uuid](https://github.com/google/uuid)
- **Logging**: Built-in `slog` (Go 1.21+)
- **Validation**: Custom validation logic

## Error Handling

The API returns consistent JSON error responses with the following structure:
```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable error message"
  }
}
```

Error codes include:
- `invalid_id`: Invalid UUID format
- `invalid_body`: Malformed request JSON
- `validation_error`: Input validation failed
- `not_found`: Resource not found
- `internal_error`: Unexpected server error

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with Go's excellent standard library
- Thanks to the creators of all open-source dependencies used