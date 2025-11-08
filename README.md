# Todo API Project

A RESTful Todo API built with Go, featuring JWT authentication, rate limiting, and MySQL database integration.

## Features

- CRUD operations for todos with JWT authentication
- Rate limiting (5 req/sec) to prevent abuse
- MySQL database with GORM ORM
- Docker support with multi-stage builds
- Health checks for Kubernetes (liveness & readiness probes)
- Graceful shutdown and CORS support
- Transaction logging with Transaction-Id header

## Tech Stack

- **Language:** Go 1.24
- **Framework:** Gin, GORM
- **Database:** MySQL/MariaDB
- **Auth:** JWT with 5-minute expiry
- **Rate Limiting:** golang.org/x/time/rate

## Quick Start

### Prerequisites

- Go 1.24+
- MySQL/MariaDB
- Docker (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/GibGyb/todo-project.git
cd todo-project

# Install dependencies
go mod download
```

### Configuration

Required environment variables:

| Variable | Description |
|----------|-------------|
| `PORT` | Server port |
| `SIGN` | JWT signing secret |
| `DB_CONN` | MySQL connection string |

**Note:** Contact the project maintainer for configuration values.

### Running

**Local Development:**
```bash
make maria          # Start MariaDB in Docker
go run main.go      # Run application (http://localhost:8081)
```

**Using Docker:**
```bash
make image          # Build Docker image
make container      # Run container
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/ping` | Health check |
| GET | `/healthz` | Readiness probe |
| GET | `/limitz` | Rate limit test |
| GET | `/x` | Build info |
| GET | `/tokenz` | Generate JWT token |

### Protected Endpoints (Requires JWT)

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/todos` | Create todo | `{"text": "Your todo"}` |
| GET | `/todos` | List all todos | - |
| DELETE | `/todos/:id` | Delete todo | - |

## Usage Example

**1. Get access token:**
```bash
curl http://localhost:8081/tokenz
```

**2. Create a todo:**
```bash
curl -X POST http://localhost:8081/todos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"text": "Buy groceries"}'
```

**3. List todos:**
```bash
curl http://localhost:8081/todos \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Project Structure

```
.
├── auth/                  # Authentication package
│   ├── auth.go           # JWT token generation
│   └── protect.go        # JWT middleware for protected routes
├── todo/                 # Todo package
│   └── todo.go          # Todo handlers and models
├── test/                # HTTP test files
│   ├── token.http       # Token generation test
│   ├── create_new_todo.http
│   ├── get_list.http
│   └── delete_todo.http
├── Dockerfile           # Multi-stage Docker build
├── Makefile            # Build and deployment commands
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Key Features

**Authentication:** JWT tokens with 5-minute expiry and `GibGyb` audience claim.

**Rate Limiting:** 5 requests/second with burst capacity of 5. Returns 429 on exceed.

**Health Checks:** Liveness (`/tmp/live`) and readiness (`/healthz`) probes for K8s.

**Graceful Shutdown:** Handles SIGINT/SIGTERM with 5-second grace period.

**CORS:** Allowed origin `http://localhost:8080` with custom headers support.

## Notes

- Todo with text "sleep" is blocked (business rule)
- Soft deletes enabled (GORM `DeletedAt`)
- Automatic database migrations on startup
- HTTP test files available in `test/` directory


## Author

GibGyb
