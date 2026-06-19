## 1. Discovery and Setup

- [ ] 1.1 Inspect the existing DB server/schema from earlier submissions and record the Krankenhaus tables, columns, IDs, and version fields needed by the Go service.
- [ ] 1.2 Initialize the Go module and choose the module path used by the repository.
- [ ] 1.3 Add required Go dependencies for `chi`, PostgreSQL access, validation support if used, and test support.
- [ ] 1.4 Create the planned package structure under `cmd/server` and `internal`.

## 2. Server Foundation

- [ ] 2.1 Implement the app assembly function with global middleware and router mounting.
- [ ] 2.2 Implement the executable server entry point with configurable bind address and graceful startup errors.
- [ ] 2.3 Add a simple health endpoint.
- [ ] 2.4 Implement a Problem Details helper for client errors with `application/problem+json`.

## 3. Krankenhaus REST API

- [ ] 3.1 Define Krankenhaus DTOs/models for REST input and output.
- [ ] 3.2 Implement `GET /rest/krankenhaus/{id}` with `ETag` and `If-None-Match` handling.
- [ ] 3.3 Implement `GET /rest/krankenhaus` for list/page behavior and `count-only`.
- [ ] 3.4 Implement `POST /rest/krankenhaus` returning HTTP 201 and `Location`.
- [ ] 3.5 Ensure HTTP handlers delegate business behavior to the service layer.

## 4. Validation

- [ ] 4.1 Implement create-payload validation for required fields, string length, email, non-negative counts, and address postal code.
- [ ] 4.2 Decide and document whether unknown JSON fields are rejected or ignored.
- [ ] 4.3 Map validation failures to HTTP 422 Problem Details responses.

## 5. PostgreSQL Persistence

- [ ] 5.1 Implement configuration loading for PostgreSQL connection settings.
- [ ] 5.2 Create a PostgreSQL connection pool during app startup.
- [ ] 5.3 Implement repository methods for finding by ID, listing/counting, and creating hospitals.
- [ ] 5.4 Map missing rows and persistence errors to service/domain errors.
- [ ] 5.5 Keep SQL and database access out of HTTP handlers.

## 6. Integration Tests

- [ ] 6.1 Add a test app setup that can serve requests without a manually started external server process.
- [ ] 6.2 Add an integration test for successful hospital creation.
- [ ] 6.3 Add an integration test for reading an existing hospital.
- [ ] 6.4 Add an integration test for invalid create input returning HTTP 422 Problem Details.
- [ ] 6.5 Run `go test ./...` and fix failures.

## 7. README and Optional Follow-up

- [ ] 7.1 Fill the README sections for REST framework, validation, PostgreSQL access, integration test, KI tools, and prompts used.
- [ ] 7.2 Document Keycloak/OIDC as optional and not required for the first mandatory server slice.
- [ ] 7.3 If time remains and the user asks for it, add an optional follow-up plan for OIDC middleware.
