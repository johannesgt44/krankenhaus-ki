## 1. Discovery and Setup

- [x] 1.1 Record the supplied PostgreSQL schema for `krankenhaus`, `adresse`, and `fachbereich`, including identity IDs, version, timestamps, indexes, checks, and cascade relationships.
- [x] 1.2 Initialize the Go module and choose the module path used by the repository.
- [x] 1.3 Add required Go dependencies for `chi`, GORM, the GORM PostgreSQL driver, validation support if used, and test support.
- [x] 1.4 Create the planned Go-typical package structure under `cmd/server` and `internal`, using German Krankenhaus domain package/type naming where appropriate.

## 2. Server Foundation

- [x] 2.1 Implement the app assembly function with global middleware and router mounting.
- [x] 2.2 Implement the executable server entry point with configurable bind address and graceful startup errors.
- [x] 2.3 Add a simple health endpoint.
- [x] 2.4 Implement a Problem Details helper for client errors with `application/problem+json`.

## 3. Krankenhaus REST API

- [x] 3.1 Define German Krankenhaus DTOs/models for REST input and output, including fields such as `Adresse`, `Fachbereich`, `Mitarbeiteranzahl`, and `Bettenanzahl`.
- [x] 3.2 Implement `GET /rest/krankenhaus/{id}` with `ETag` and `If-None-Match` handling.
- [x] 3.3 Implement `GET /rest/krankenhaus` for list/page behavior and `count-only`.
- [x] 3.4 Implement `POST /rest/krankenhaus` returning HTTP 201 and `Location`.
- [x] 3.5 Implement `PUT /rest/krankenhaus/{id}` with `If-Match`, version handling, HTTP 204, and `ETag`.
- [x] 3.6 Implement `DELETE /rest/krankenhaus/{id}` returning HTTP 204.
- [x] 3.7 Ensure HTTP handlers delegate business behavior to the service layer.

## 4. Validation

- [x] 4.1 Implement create/update payload validation for German JSON fields, required fields, string length, email, non-negative counts, and address postal code.
- [x] 4.2 Decide and document whether unknown JSON fields are rejected or ignored.
- [x] 4.3 Map validation failures to HTTP 422 Problem Details responses.

## 5. PostgreSQL Persistence

- [x] 5.1 Implement configuration loading for the Docker PostgreSQL connection settings, including host, port `5432`, user `postgres`, and password from environment/config.
- [x] 5.2 Create a GORM PostgreSQL connection during app startup.
- [x] 5.3 Implement explicit GORM models for `krankenhaus`, `adresse`, and `fachbereich` mapped to the supplied schema.
- [x] 5.4 Add SQL reset and seed scripts for development/test startup using the supplied schema.
- [x] 5.5 Add a configuration flag that enables database reset/seed at startup for development/test and leaves existing tables untouched otherwise.
- [x] 5.6 Implement repository methods for finding by ID, listing/counting, creating, updating, and deleting hospitals.
- [x] 5.7 Map missing rows and persistence errors to service/domain errors.
- [x] 5.8 Keep ORM/database access out of HTTP handlers.

## 6. Integration Tests

- [x] 6.1 Add a test app setup that can serve requests without a manually started external server process.
- [x] 6.2 Add an integration test for successful hospital creation.
- [x] 6.3 Add an integration test for reading an existing hospital.
- [x] 6.4 Add an integration test for invalid create input returning HTTP 422 Problem Details.
- [x] 6.5 Add an integration test for updating a hospital with `If-Match`.
- [x] 6.6 Add an integration test for deleting a hospital.
- [x] 6.7 Run `go test ./...` and fix failures.

## 7. README and Optional Follow-up

- [x] 7.1 Fill the README sections in German for REST framework, Go-typical structure, German domain naming, validation, GORM/PostgreSQL access, integration test, KI tools, and prompts used.
- [x] 7.2 Document Keycloak/OIDC as optional and not required for the first mandatory server slice.
- [x] 7.3 If time remains and the user asks for it, add an optional follow-up plan for OIDC middleware.
