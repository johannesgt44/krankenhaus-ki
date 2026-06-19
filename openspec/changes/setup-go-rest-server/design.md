## Context

The repository currently has README placeholders and OpenSpec configuration, but no Go server code. The assignment attachment sets Go as the platform, expects use of the existing DB server from earlier submissions, and marks Keycloak as optional. The README asks for framework/library choices, REST read/create functionality, validation for create, PostgreSQL ORM/persistence, optional OIDC with Keycloak, a simple integration test, and recorded AI prompts.

The REST reference map points to TypeScript/Hono and Python/FastAPI implementations. The transferable patterns are: app entry point registers global middleware and feature routers, read/write routes are separated, handlers call services, request models validate JSON, client errors use Problem Details, `GET` supports `ETag`/`If-None-Match`, `POST` returns `201 Created` with `Location`, and integration tests cover REST behavior.

## Goals / Non-Goals

**Goals:**
- Create a small Go REST server plan that can be implemented after approval.
- Use an idiomatic Go router/framework choice for maintainable REST routes.
- Cover the assignment's mandatory scope: read hospitals, create hospitals, validate create input, connect PostgreSQL, and add a simple integration test.
- Keep the architecture close to the reference-map patterns while adapting them to Go.

**Non-Goals:**
- Do not implement code until the user says `apply`.
- Do not require Keycloak/OIDC for the first runnable server.
- Do not include GraphQL, Prometheus, shutdown/admin endpoints, or full production observability in the first slice.
- Do not require full update/delete behavior unless time remains after mandatory read/create behavior.

## Decisions

### Use Go with `chi` plus `net/http` for REST routing

Choose `github.com/go-chi/chi/v5` as the REST router/framework layer. `chi` stays compatible with `net/http`, has middleware and route grouping, supports modular routers, and maps well to the Hono/FastAPI reference structure without a large framework surface.

Alternatives considered:
- Standard `net/http.ServeMux`: very attractive after Go 1.22 because method/path patterns and wildcards are built in. Use it if dependency minimization becomes the priority. It is less convenient for grouped routers and middleware composition than `chi`.
- Gin or Echo: productive and popular, but bring more framework conventions than needed for a workshop server.
- Fiber: fast and ergonomic, but it is not built on `net/http`, which makes standard middleware/testing reuse less direct.

### Use a layered package layout

Use packages that mirror the reference pattern:

```text
cmd/server              # main entry point
internal/app            # app assembly, middleware, router mounting
internal/hospital/http  # read/create handlers and DTOs
internal/hospital/service
internal/hospital/repository
internal/problem        # Problem Details helpers
```

Handlers parse HTTP concerns and call services. Services enforce business flow. Repositories isolate PostgreSQL details. This keeps the implementation explainable for the seminar and easy to test.

### Separate read and write route registration

Mount read and write routes under `/rest/krankenhaus`, but keep their registration separate in code, matching the reference-map split between read routers and write routers. The first mandatory write route is `POST /rest/krankenhaus`; update/delete can remain a later extension.

### Return Problem Details for client errors

Use `application/problem+json` for validation failures, not found, precondition failures, and malformed request bodies. Use a small local helper instead of a broad dependency.

### Use PostgreSQL through `pgx`

Use `github.com/jackc/pgx/v5` or its pool package for PostgreSQL connectivity. It is a direct, well-supported Go PostgreSQL driver and avoids introducing a heavy ORM before the DB shape is confirmed. If the assignment strictly requires an ORM label in the README, note this as "PostgreSQL access via pgx repository layer" or add a lightweight query helper only if needed.

### Keep Keycloak optional

Design the routing so auth middleware can be inserted later, but do not block the mandatory REST server on OIDC. If implemented, OIDC should be a separate follow-up capability using middleware around write routes or protected reads.

## Risks / Trade-offs

- DB schema mismatch -> Mitigation: inspect the previous DB server/schema before implementing repository SQL.
- Assignment expects an "ORM" by name -> Mitigation: confirm whether pgx repository access is accepted; otherwise add a small Go ORM/query layer during apply.
- Overbuilding beyond the workshop -> Mitigation: keep first implementation to read/create/validation/PostgreSQL/integration test.
- Keycloak time sink -> Mitigation: keep optional OIDC out of the first apply slice unless explicitly requested.
- REST edge cases around ETag/versioning -> Mitigation: implement ETag for `GET /{id}` and design version handling so `PUT` can be added later without route redesign.
