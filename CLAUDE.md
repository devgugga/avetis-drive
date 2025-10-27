# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based application called "avetis-drive" following Clean Architecture principles with a clear separation of concerns across domain, use cases, infrastructure, and presentation layers.

## Development Commands

### Building and Running

- Build the application: `go build -o ./tmp/main.exe ./cmd/api`
- Run with hot reload (Air): `air`
  - Air is configured via `.air.toml`
  - Watches `.go`, `.tpl`, `.tmpl`, `.html` files
  - Excludes test files and tmp/vendor/testdata directories
  - Output binary: `tmp/main.exe`
  - Build errors logged to: `build-errors.log`

### Testing

- Run all tests: `go test ./...`
- Run tests with coverage: `go test -cover ./...`
- Run tests for a specific package: `go test ./internal/domain/...`

### Dependency Management

- Install/update dependencies: `go mod tidy`
- Vendor dependencies: `go mod vendor`
- View dependencies: `go list -m all`

### Database (Ent + PostgreSQL)

- Generate Ent code after schema changes: `go generate ./internal/infrastructure/database/ent`
- Create new Ent schema: `go run entgo.io/ent/cmd/ent init --target internal/infrastructure/database/ent/schema <EntityName>`
- Database migrations are run automatically on application startup via `AutoMigrate()`
- **Database auto-creation**: The application automatically creates the PostgreSQL database if it doesn't exist (requires PostgreSQL server to be running)

## Architecture

The codebase follows **Clean Architecture** with clear boundaries between layers:

### Directory Structure

```
cmd/api/          - Application entry point (main.go)
internal/
  domain/         - Core business entities and logic (innermost layer)
  usecases/       - Application business rules and orchestration
  dto/            - Data Transfer Objects
    requests/     - API request DTOs
    responses/    - API response DTOs
  infrastructure/ - External concerns (outermost layer)
    config/       - Configuration management (.env support via godotenv)
    logging/      - Structured logging (zerolog)
    database/     - Database layer (Ent ORM + PostgreSQL)
      ent/        - Generated Ent code
        schema/   - Ent schema definitions
    http/         - HTTP layer (Echo framework)
      handlers/   - HTTP request handlers
      middlewares/- HTTP middlewares
```

### Architectural Principles

1. **Dependency Rule**: Dependencies point inward. Domain layer has no dependencies on other layers. Use cases depend only on domain. Infrastructure depends on use cases and domain.

2. **Domain Layer** (`internal/domain/`): Contains enterprise business rules, entities, and domain services. Should be framework-agnostic and have no external dependencies.

3. **Use Cases Layer** (`internal/usecases/`): Contains application-specific business rules. Orchestrates data flow between domain entities and coordinates domain logic.

4. **Infrastructure Layer** (`internal/infrastructure/`): Handles external concerns like HTTP, databases, and configuration. Implements interfaces defined in inner layers.

5. **DTO Pattern**: Request/response objects separate API contracts from domain models, preventing external API changes from affecting core business logic.

### Adding New Features

When adding new features, follow this order:

1. Define domain entities in `internal/domain/`
2. Create use case interfaces and implementations in `internal/usecases/`
3. Define request/response DTOs in `internal/dto/`
4. Implement HTTP handlers in `internal/infrastructure/http/handlers/`
5. Wire up dependencies in `cmd/api/main.go`

## Module Information

- Module path: `github.com/devgugga/avetis-drive`
- Go version: 1.25
