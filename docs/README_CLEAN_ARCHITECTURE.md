# Clean Architecture Implementation

AntiquaOS Backend strictly follows Clean Architecture structure to keep the code heavily decoupled and testable.

## Layers structure

1. **Domain Layer** (`src/domain/`)
   Contains central business structures and models (`Product`, `Category`, `Sale`, etc.). Has no external dependencies.
2. **Application Layer** (`src/application/usecases/`)
   Contains business logic orchestration (Use Cases). Interfaces with repositories but does not know their SQL implementation.
3. **Infrastructure Layer** (`src/infrastructure/`)
   Contains the technical implementation details:
   - `rest/controllers`: Gin HTTP handlers mapping payloads to domain models.
   - `repository/psql`: Concrete GORM/PostgreSQL adapters matching UserUseCase interfaces.
   - `di`: Global Dependency Injection wireup (`application_context.go`).
   - `websocket`: Gorilla WS Hub and Clients for real-time data streaming.
   - `security`: JWT and password crypto hashing.
   - `logger`: Uber Zap log formats.