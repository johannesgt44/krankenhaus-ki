## ADDED Requirements

### Requirement: Server starts from a Go entry point
The system SHALL provide a Go server entry point that starts an HTTP server on a configurable address.

#### Scenario: Start server with default configuration
- **WHEN** the server is started without an explicit bind address
- **THEN** it listens on the documented default address

#### Scenario: Start server with configured address
- **WHEN** a bind address is supplied through configuration
- **THEN** the server listens on that configured address

### Requirement: Server mounts REST routers under a stable prefix
The system SHALL mount Krankenhaus REST routes under `/rest/krankenhaus`.

#### Scenario: REST prefix is available
- **WHEN** a request targets `/rest/krankenhaus`
- **THEN** the request is routed to the Krankenhaus REST handlers

### Requirement: Server provides a health endpoint
The system SHALL expose a simple health endpoint for checking that the process is running.

#### Scenario: Health check succeeds
- **WHEN** a client sends a request to the health endpoint
- **THEN** the server responds with HTTP 200

### Requirement: Server normalizes client errors
The system SHALL return client-facing errors as `application/problem+json` Problem Details responses.

#### Scenario: Not found uses Problem Details
- **WHEN** a REST request addresses a missing resource
- **THEN** the server responds with HTTP 404 and a Problem Details JSON body

#### Scenario: Invalid request uses Problem Details
- **WHEN** a REST request contains invalid JSON or invalid fields
- **THEN** the server responds with HTTP 422 and a Problem Details JSON body
