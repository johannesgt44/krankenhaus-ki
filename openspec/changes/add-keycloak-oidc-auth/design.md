## Context

The current Go REST server exposes all hospital endpoints without authentication. Routing is assembled in `internal/app`, the hospital REST handlers live under `internal/krankenhaus/rest`, and configuration is loaded through `internal/config`.

The reference project already contains a Keycloak Docker Compose setup under `extras/compose/keycloak`. This Go project should follow the same local-development idea, but keep a Go-typical structure with authentication implemented as middleware rather than embedding security checks inside every handler.

## Goals / Non-Goals

**Goals:**

- Validate Keycloak OIDC access tokens for protected hospital endpoints.
- Require an `admin` role for create, update, and delete operations.
- Keep `GET /health`, `GET /rest/krankenhaus`, and `GET /rest/krankenhaus/{id}` public.
- Add local Keycloak container support under `extras/compose`.
- Configure OIDC through environment variables.
- Keep REST handler code focused on HTTP/domain behavior.

**Non-Goals:**

- No user registration or user-management API in the Go server.
- No token issuing endpoint in the Go server.
- No frontend login flow.
- No production-grade Keycloak realm export/import automation unless needed during implementation.
- No fine-grained distinction between create, update, and delete roles in the first version.

## Decisions

### Use middleware for authentication and authorization

Protected routes will be wrapped with middleware from a new package such as `internal/security`. The middleware will:

- read the `Authorization: Bearer <token>` header.
- validate the JWT issuer, signature, expiry, and configured audience/client.
- extract roles from Keycloak claims.
- require the `admin` role for protected routes.

This keeps `internal/krankenhaus/rest` free of Keycloak-specific code. The main alternative was checking roles inside each handler, but that would duplicate security behavior and mix infrastructure concerns into REST endpoint logic.

### Protect only mutating hospital routes

`internal/app` or the mounted hospital router will explicitly separate public and protected routes:

- public: `GET /health`, `GET /rest/krankenhaus`, `GET /rest/krankenhaus/{id}`
- protected with `admin`: `POST /rest/krankenhaus`, `PUT /rest/krankenhaus/{id}`, `DELETE /rest/krankenhaus/{id}`

The alternative was protecting the complete `/rest/krankenhaus` mount and adding exceptions for search. Explicit route grouping is easier to test and avoids accidentally locking public search.

### Use Keycloak client role `admin`

The planned local Keycloak setup uses:

- realm: `python`
- client: `python-client`
- role: `admin`

The middleware should prefer client roles from `resource_access["python-client"].roles`. Accepting `realm_access.roles` as a fallback is acceptable if it simplifies manual local setup, but tests should cover the configured client-role behavior.

### Add OIDC configuration under `internal/config`

Configuration will be loaded from environment variables, for example:

- `OIDC_ENABLED`
- `OIDC_ISSUER_URL`
- `OIDC_CLIENT_ID`
- `OIDC_REQUIRED_ROLE`

Recommended defaults for local development:

- `OIDC_ENABLED=true`
- `OIDC_ISSUER_URL=http://localhost:8880/realms/python`
- `OIDC_CLIENT_ID=python-client`
- `OIDC_REQUIRED_ROLE=admin`

Keeping `OIDC_ENABLED=true` by default makes protected routes active for normal workshop starts. Developers can set `OIDC_ENABLED=false` explicitly when they want to run the server without local Keycloak.

### Use a standard OIDC/JWT validation library

Use a maintained Go OIDC/JWT library rather than handwritten JWT parsing. The implementation should verify tokens cryptographically through the Keycloak JWKS endpoint discovered from the issuer metadata.

The main alternative is manual JWKS loading and claim verification. That increases security risk and maintenance cost.

### Add Keycloak to local Compose support

Extend `extras/compose` with Keycloak support oriented on the reference project:

- Keycloak service published on local HTTP port `8880`.
- bootstrap admin user `tmp` with password `p` for development.
- persistent Keycloak data volume.
- PostgreSQL-backed Keycloak storage if practical with the existing Compose setup.

If PostgreSQL setup becomes too heavy for the workshop scope, the implementation may use Keycloak development storage initially, but the docs must state that it is for local development only.

## Risks / Trade-offs

- Keycloak claim shape can differ depending on realm/client-role configuration -> document the expected role mapping and test the claim parser with representative claims.
- Keeping `OIDC_ENABLED=true` by default requires local Keycloak for normal protected-route startup, but avoids accidentally testing write endpoints without auth -> README must clearly show how to disable OIDC for local troubleshooting.
- Keycloak with PostgreSQL adds Compose complexity -> keep commands and defaults close to the reference project and document manual realm/client setup.
- Local issuer URLs differ between host and Docker networking -> support configurable issuer URL and document the host-default value.

## Migration Plan

1. Add OIDC configuration fields and defaults.
2. Add security middleware with token validation and role checks.
3. Wire public/protected hospital routes.
4. Add tests for public search, missing token, invalid role, and admin role success.
5. Add Keycloak Compose support and setup documentation.
6. Update README with OIDC startup and token usage.

Rollback is straightforward: disable OIDC through `OIDC_ENABLED=false` or remove middleware wiring while leaving the rest of the REST API unchanged.

## Open Questions

- Should local Keycloak realm/client/users be created manually through the UI, imported from a realm JSON file, or scripted later?
- Should `admin` be accepted only as a client role, or also as a realm role for easier manual setup?
