## Why

The REST API currently allows create, update, and delete operations without authentication. OIDC with Keycloak should protect mutating operations while keeping hospital search endpoints public for read access and workshop usability.

## What Changes

- Add OIDC bearer-token authentication using Keycloak as the local authorization server.
- Add role-based authorization for protected hospital write operations.
- Keep search and health endpoints public:
  - `GET /health`
  - `GET /rest/krankenhaus`
  - `GET /rest/krankenhaus/{id}`
- Require a valid Keycloak access token with the `admin` role for:
  - `POST /rest/krankenhaus`
  - `PUT /rest/krankenhaus/{id}`
  - `DELETE /rest/krankenhaus/{id}`
- Add local Keycloak Docker Compose support under `extras/compose`, oriented on the existing reference project setup.
- Add configuration through environment variables so OIDC can be enabled and configured without code changes.

## Capabilities

### New Capabilities
- `keycloak-oidc-auth`: OIDC bearer-token validation and role-based protection for selected REST endpoints.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `cmd/server`
  - `internal/app`
  - new authentication/authorization package under `internal`
  - configuration under `internal/config`
  - REST tests under `internal/app` or package-local tests
- Affected API behavior:
  - unauthenticated create, update, and delete requests will return an authorization error.
  - search and health endpoints remain publicly accessible.
- Affected dependencies:
  - OIDC/JWT validation dependency for Go.
  - Keycloak Docker Compose service and local setup documentation.
- Affected external systems:
  - local Keycloak container with a `krankenhaus` realm, `krankenhaus-client` client, and admin role assignments.
