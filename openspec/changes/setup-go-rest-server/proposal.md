## Why

The project needs a new Go-based REST server for the programming workshop assignment. The current repository only contains planning/configuration material and the README placeholders, so the first useful step is to define a small, runnable server plan before implementation starts.

## What Changes

- Add a Go HTTP server for the Krankenhaus domain.
- Use a REST framework/router suitable for idiomatic Go services.
- Use a Go-typical project structure while keeping German Krankenhaus domain terminology in types, DTOs, JSON fields, and user-facing error details.
- Provide REST endpoints for reading, creating, updating, and deleting hospitals.
- Validate JSON input for creating and updating hospitals.
- Connect the service to the existing Docker-based PostgreSQL DB server using GORM.
- Normalize REST errors as Problem Details JSON.
- Add a simple integration test covering the required REST path.
- Keep OIDC/Keycloak optional and out of the mandatory first implementation slice.

## Capabilities

### New Capabilities
- `go-rest-server`: Runnable Go REST service entry point, middleware, health behavior, routing, and error response conventions.
- `hospital-rest-api`: Krankenhaus REST resources for read, create, update, and delete flows, including status codes, headers, and JSON shapes.
- `hospital-validation`: Request validation rules for creating and updating hospitals.
- `postgres-persistence`: PostgreSQL-backed persistence access through GORM and a service/repository boundary.
- `rest-integration-test`: Minimal integration checks for the REST server and required endpoints.

### Modified Capabilities
- None.

## Impact

- New Go module, server entry point, router, handlers, service layer, repository layer, DTO/model definitions, validation helpers, Problem Details helper, and tests.
- Implementation naming should use German domain terms such as `Krankenhaus`, `Adresse`, `Fachbereich`, `Bettenanzahl`, `Mitarbeiteranzahl`, and ASCII JSON field names such as `strasse`.
- New dependencies are expected for HTTP routing, GORM/PostgreSQL access, and test assertions or HTTP testing support.
- README will later need to be filled with selected frameworks/libraries, KI tool usage, repository link, and prompts, but that belongs to implementation/documentation tasks after this proposal is accepted.
