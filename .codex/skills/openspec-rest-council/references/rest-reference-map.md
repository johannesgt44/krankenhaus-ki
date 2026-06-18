# REST Reference Map

Use these repositories as reference implementations when preparing REST server changes with OpenSpec.

## TypeScript/Hono

Repository: <https://github.com/johannesgt44/krankenhaus>

Look for:

- `src/app.mts` for app setup, middleware, router registration, route mounting, and error handling.
- `src/krankenhaus/router/krankenhaus-router.mts` for read routes.
- `src/krankenhaus/router/krankenhaus-write-router.mts` for create, update, and delete routes.
- `src/krankenhaus/router/krankenhaus-validation.mts` for request validation with Zod.
- `src/krankenhaus/service/` for business logic.
- `src/problem-details.mts` for structured HTTP error responses.
- `test/integration/rest/` for integration test examples.

## Python/FastAPI

Repository: <https://github.com/marculm/krankenhaus>

Look for:

- `src/krankenhaus/fastapi_app.py` for app setup, middleware, router registration, and exception handlers.
- `src/krankenhaus/router/krankenhaus_router.py` for read routes.
- `src/krankenhaus/router/krankenhaus_write_router.py` for create, update, and delete routes.
- `src/krankenhaus/router/*_model.py` for Pydantic request models.
- `src/krankenhaus/router/dependencies.py` for dependency injection.
- `src/krankenhaus/service/` for business logic.
- `src/krankenhaus/problem_details.py` for structured HTTP error responses.
- `tests/integration/rest/` for integration test examples.

## Shared Patterns

Extract these patterns before writing an OpenSpec proposal:

- App entry point registers global middleware and feature routers.
- Read and write REST routers are separated.
- Routers call services instead of directly containing business logic.
- Request models or schemas validate incoming JSON.
- Errors are normalized as Problem Details.
- `GET /resource/{id}` handles not found and may use `ETag` / `If-None-Match`.
- `POST /resource` returns `201 Created` and a `Location` header.
- `PUT /resource/{id}` uses `If-Match` when optimistic locking is required.
- `DELETE /resource/{id}` returns `204 No Content`.
- Integration tests cover GET, POST, PUT, DELETE, invalid input, and not-found cases.
