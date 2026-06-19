## Why

The project needs a new Go-based REST server for the programming workshop assignment. The current repository only contains planning/configuration material and the README placeholders, so the first useful step is to define a small, runnable server plan before implementation starts.

## What Changes

- Add a Go HTTP server for the Krankenhaus domain.
- Use a REST framework/router suitable for idiomatic Go services.
- Provide REST endpoints for reading hospitals and creating a new hospital.
- Validate JSON input for creating hospitals.
- Connect the service to the existing PostgreSQL DB server from earlier submissions.
- Normalize REST errors as Problem Details JSON.
- Add a simple integration test covering the required REST path.
- Keep OIDC/Keycloak optional and out of the mandatory first implementation slice.

## Capabilities

### New Capabilities
- `go-rest-server`: Runnable Go REST service entry point, middleware, health behavior, routing, and error response conventions.
- `hospital-rest-api`: Krankenhaus REST resources for read and create flows, including status codes, headers, and JSON shapes.
- `hospital-validation`: Request validation rules for creating hospitals.
- `postgres-persistence`: PostgreSQL-backed persistence access through a service/repository boundary.
- `rest-integration-test`: Minimal integration checks for the REST server and required endpoints.

### Modified Capabilities
- None.

## Impact

- New Go module, server entry point, router, handlers, service layer, repository layer, DTO/model definitions, validation helpers, Problem Details helper, and tests.
- New dependencies are expected for HTTP routing, PostgreSQL access, and test assertions or HTTP testing support.
- README will later need to be filled with selected frameworks/libraries, KI tool usage, repository link, and prompts, but that belongs to implementation/documentation tasks after this proposal is accepted.
