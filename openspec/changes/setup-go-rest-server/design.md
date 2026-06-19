## Context

The repository currently has README placeholders and OpenSpec configuration, but no Go server code. The assignment attachment sets Go as the platform, expects use of the existing DB server from earlier submissions, and marks Keycloak as optional. The README asks for framework/library choices, REST read/create functionality, validation for create, PostgreSQL ORM/persistence, optional OIDC with Keycloak, a simple integration test, and recorded AI prompts. The user has supplied the PostgreSQL schema, Docker database image/port, credentials, and asked for GORM, update/delete support, and German wording throughout.

The REST reference map points to TypeScript/Hono and Python/FastAPI implementations. The transferable patterns are: app entry point registers global middleware and feature routers, read/write routes are separated, handlers call services, request models validate JSON, client errors use Problem Details, `GET` supports `ETag`/`If-None-Match`, `POST` returns `201 Created` with `Location`, `PUT` uses `If-Match`, `DELETE` returns `204 No Content`, and integration tests cover REST behavior.

## Goals / Non-Goals

**Goals:**
- Create a small Go REST server plan that can be implemented after approval.
- Use an idiomatic Go router/framework choice for maintainable REST routes.
- Use Go-typical package boundaries and German domain naming for the Krankenhaus API.
- Cover read, create, update, delete, validation, PostgreSQL via GORM, and a simple integration test.
- Keep the architecture close to the reference-map patterns while adapting them to Go.

**Non-Goals:**
- Do not implement code until the user says `apply`.
- Do not require Keycloak/OIDC for the first runnable server.
- Do not include GraphQL, Prometheus, shutdown/admin endpoints, or full production observability in the first slice.
- Do not include GraphQL, Prometheus, shutdown/admin endpoints, or full production observability in the first slice.

## Decisions

### Use Go with `chi` plus `net/http` for REST routing

Choose `github.com/go-chi/chi/v5` as the REST router/framework layer. `chi` stays compatible with `net/http`, has middleware and route grouping, supports modular routers, and maps well to the Hono/FastAPI reference structure without a large framework surface.

Alternatives considered:
- Standard `net/http.ServeMux`: very attractive after Go 1.22 because method/path patterns and wildcards are built in. Use it if dependency minimization becomes the priority. It is less convenient for grouped routers and middleware composition than `chi`.
- Gin or Echo: productive and popular, but bring more framework conventions than needed for a workshop server.
- Fiber: fast and ergonomic, but it is not built on `net/http`, which makes standard middleware/testing reuse less direct.

### Use a Go-typical layered package layout

Use packages that mirror the reference pattern while staying idiomatic for Go: small packages, explicit constructors, interfaces at the consumer side when useful, and no framework-specific business logic. Package names stay lowercase and ASCII. Domain names in exported types, DTOs, methods, and JSON tags use German terminology.

```text
cmd/server                  # main entry point
internal/app                # app assembly, middleware, router mounting
internal/konfiguration      # environment/config loading
internal/krankenhaus/rest   # read/create handlers and DTOs
internal/krankenhaus/service
internal/krankenhaus/repository
internal/problem            # Problem Details helpers
```

Handlers parse HTTP concerns and call services. Services enforce business flow. Repositories isolate PostgreSQL details. This keeps the implementation explainable for the seminar and easy to test.

### Use German domain terminology in code and JSON

Use German names for the Krankenhaus domain, for example `Krankenhaus`, `Adresse`, `Fachbereich`, `KrankenhausService`, `KrankenhausRepository`, `SucheNachID`, `Suche`, and `Erstellen`. JSON fields should match the German reference shape: `name`, `mitarbeiteranzahl`, `bettenanzahl`, `email`, `adresse`, `strasse`, `hausnummer`, `plz`, `ort`, and `fachbereiche`.

Go keywords, standard-library APIs, dependency APIs, package conventions, and test function prefixes remain normal Go. Avoid umlauts in code and JSON identifiers; use ASCII spellings like `strasse`.

### Separate read and write route registration

Mount read and write routes under `/rest/krankenhaus`, but keep their registration separate in code, matching the reference-map split between read routers and write routers. Write routes include `POST /rest/krankenhaus`, `PUT /rest/krankenhaus/{id}`, and `DELETE /rest/krankenhaus/{id}`.

### Return Problem Details for client errors

Use `application/problem+json` for validation failures, not found, precondition failures, and malformed request bodies. Use a small local helper instead of a broad dependency.

### Use PostgreSQL through GORM

Use `gorm.io/gorm` with `gorm.io/driver/postgres` for ORM mapping because the assignment asks for ORM and the user explicitly chose GORM. Map the existing schema instead of letting GORM freely auto-migrate it. Models should use explicit table names and column tags for `krankenhaus`, `adresse`, and `fachbereich`.

The known database uses `krankenhausspace` as default tablespace and these tables:

```text
krankenhaus(id, version, name, mitarbeiteranzahl, bettenanzahl, email, erzeugt, aktualisiert)
adresse(id, strasse, hausnummer, plz, ort, krankenhaus_id)
fachbereich(id, name, beschreibung, leitung, anzahlaerzte, krankenhaus_id)
```

The Docker database is reachable on published port `5432` with database/user/password values to be handled as configuration. The supplied credentials are user `postgres` and password `p`; do not hard-code the password in source code.

### Use German API messages and documentation

All user-facing error details, README content, and API examples should be German. Code identifiers use German domain terms where sensible, but Go syntax and common Go naming rules remain idiomatic.

### Keep Keycloak optional

Design the routing so auth middleware can be inserted later, but do not block the mandatory REST server on OIDC. If implemented, OIDC should be a separate follow-up capability using middleware around write routes or protected reads.

## Risks / Trade-offs

- DB schema mismatch -> Mitigation: use the supplied schema as implementation source and verify against the running Docker database during apply.
- German naming can become awkward in generic technical packages -> Mitigation: keep technical package names short/lowercase and use German naming mainly for domain types, DTOs, service methods, JSON fields, and error details.
- GORM auto-migration could drift from the existing schema -> Mitigation: prefer explicit models and avoid destructive auto-migration; create schema only if a local test database needs initialization.
- Overbuilding beyond the workshop -> Mitigation: implement REST CRUD only, with Keycloak and extra endpoints left optional.
- Keycloak time sink -> Mitigation: keep optional OIDC out of the first apply slice unless explicitly requested.
- REST edge cases around ETag/versioning -> Mitigation: implement ETag for `GET /{id}` and design version handling so `PUT` can be added later without route redesign.
