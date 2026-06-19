## 1. Configuration

- [ ] 1.1 Add OIDC configuration fields to `internal/config` for enabled flag, issuer URL, client ID, and required role.
- [ ] 1.2 Add local defaults for Keycloak realm `krankenhaus`, client `krankenhaus-client`, and role `admin`.
- [ ] 1.3 Pass OIDC configuration from `cmd/server` into app/router setup.

## 2. Security Middleware

- [ ] 2.1 Add a Go OIDC/JWT validation dependency.
- [ ] 2.2 Create `internal/security` package for bearer-token validation and role extraction.
- [ ] 2.3 Validate issuer, signature, expiry, and configured client/audience for protected requests.
- [ ] 2.4 Extract Keycloak roles from `resource_access[client].roles` and support the configured required role.
- [ ] 2.5 Return `401 Unauthorized` for missing or invalid tokens.
- [ ] 2.6 Return `403 Forbidden` for valid tokens without the required role.

## 3. Route Protection

- [ ] 3.1 Keep `GET /health` public.
- [ ] 3.2 Keep `GET /rest/krankenhaus` public.
- [ ] 3.3 Keep `GET /rest/krankenhaus/{id}` public.
- [ ] 3.4 Protect `POST /rest/krankenhaus` with the admin role.
- [ ] 3.5 Protect `PUT /rest/krankenhaus/{id}` with the admin role.
- [ ] 3.6 Protect `DELETE /rest/krankenhaus/{id}` with the admin role.

## 4. Tests

- [ ] 4.1 Add tests showing public health and search endpoints work without a token.
- [ ] 4.2 Add tests showing protected endpoints return `401` without a token when OIDC is enabled.
- [ ] 4.3 Add tests showing protected endpoints return `403` for a valid token without the admin role.
- [ ] 4.4 Add tests showing protected endpoints accept a valid token with the admin role.
- [ ] 4.5 Add tests for Keycloak role-claim parsing.

## 5. Keycloak Compose and Documentation

- [ ] 5.1 Add Keycloak Docker Compose support under `extras/compose`, oriented on the reference project.
- [ ] 5.2 Document local Keycloak startup and access URL.
- [ ] 5.3 Document realm `krankenhaus`, client `krankenhaus-client`, role `admin`, and test users.
- [ ] 5.4 Document required environment variables for enabling OIDC in the Go server.
- [ ] 5.5 Document how to obtain a token and call protected endpoints with Bruno or curl.

## 6. Verification

- [ ] 6.1 Run `gofmt` on changed Go files.
- [ ] 6.2 Run `go test ./...`.
- [ ] 6.3 Manually verify public and protected endpoint behavior against local Keycloak if Docker is available.
