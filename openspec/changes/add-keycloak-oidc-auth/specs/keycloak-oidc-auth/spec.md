## ADDED Requirements

### Requirement: Public hospital search remains unauthenticated
The system SHALL allow public access to hospital search and health endpoints without an OIDC access token.

#### Scenario: Health endpoint without token
- **WHEN** a client requests `GET /health` without an `Authorization` header
- **THEN** the system returns `200 OK`

#### Scenario: Hospital list without token
- **WHEN** a client requests `GET /rest/krankenhaus` without an `Authorization` header
- **THEN** the system processes the search request without requiring authentication

#### Scenario: Hospital by ID without token
- **WHEN** a client requests `GET /rest/krankenhaus/{id}` without an `Authorization` header
- **THEN** the system processes the read request without requiring authentication

### Requirement: Mutating hospital endpoints require authentication
The system SHALL reject unauthenticated create, update, and delete requests for hospitals when OIDC is enabled.

#### Scenario: Create without token
- **WHEN** a client requests `POST /rest/krankenhaus` without an `Authorization` header and OIDC is enabled
- **THEN** the system returns `401 Unauthorized`

#### Scenario: Update without token
- **WHEN** a client requests `PUT /rest/krankenhaus/{id}` without an `Authorization` header and OIDC is enabled
- **THEN** the system returns `401 Unauthorized`

#### Scenario: Delete without token
- **WHEN** a client requests `DELETE /rest/krankenhaus/{id}` without an `Authorization` header and OIDC is enabled
- **THEN** the system returns `401 Unauthorized`

### Requirement: Mutating hospital endpoints require admin role
The system SHALL require a valid Keycloak access token with the configured admin role for hospital create, update, and delete operations.

#### Scenario: Valid token without admin role
- **WHEN** a client requests a mutating hospital endpoint with a valid OIDC access token that does not contain the configured admin role
- **THEN** the system returns `403 Forbidden`

#### Scenario: Valid token with admin role
- **WHEN** a client requests a mutating hospital endpoint with a valid OIDC access token that contains the configured admin role
- **THEN** the system processes the request

#### Scenario: Invalid token
- **WHEN** a client requests a mutating hospital endpoint with an invalid, expired, or unverifiable OIDC access token
- **THEN** the system returns `401 Unauthorized`

### Requirement: OIDC configuration is environment-driven
The system SHALL load OIDC settings from environment variables so local Keycloak and other OIDC issuers can be configured without code changes.

#### Scenario: OIDC disabled by default
- **WHEN** the server starts without `OIDC_ENABLED=true`
- **THEN** the system starts without requiring a reachable Keycloak issuer

#### Scenario: OIDC enabled with issuer and client
- **WHEN** the server starts with OIDC enabled and issuer/client settings configured
- **THEN** the system validates protected-route access tokens against the configured issuer and client

### Requirement: Local Keycloak support is documented
The system SHALL provide local Docker Compose support and documentation for running Keycloak during development.

#### Scenario: Developer starts local Keycloak
- **WHEN** a developer follows the documented Docker Compose instructions
- **THEN** Keycloak is reachable on the documented local port for realm and token configuration

#### Scenario: Developer configures roles
- **WHEN** a developer follows the documented Keycloak setup
- **THEN** a user can receive an access token containing the admin role required by protected REST endpoints
