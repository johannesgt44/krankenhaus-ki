## ADDED Requirements

### Requirement: Client can read a hospital by ID
The system SHALL expose `GET /rest/krankenhaus/{id}` to retrieve a hospital by numeric ID.

#### Scenario: Hospital exists
- **WHEN** a client requests an existing hospital ID
- **THEN** the server responds with HTTP 200, a hospital JSON body, and an `ETag` header

#### Scenario: Hospital does not exist
- **WHEN** a client requests a missing hospital ID
- **THEN** the server responds with HTTP 404

#### Scenario: Cached hospital version matches
- **WHEN** a client sends `If-None-Match` with the current hospital version
- **THEN** the server responds with HTTP 304 and no response body

### Requirement: Client can list hospitals
The system SHALL expose `GET /rest/krankenhaus` to search or page through hospitals.

#### Scenario: List request succeeds
- **WHEN** a client requests `/rest/krankenhaus`
- **THEN** the server responds with HTTP 200 and a JSON page or list of hospitals

#### Scenario: Count-only request succeeds
- **WHEN** a client requests `/rest/krankenhaus?count-only`
- **THEN** the server responds with HTTP 200 and a JSON object containing the total count

### Requirement: Client can create a hospital
The system SHALL expose `POST /rest/krankenhaus` to create a new hospital from a valid JSON request body.

#### Scenario: Valid hospital is created
- **WHEN** a client posts a valid hospital creation payload
- **THEN** the server responds with HTTP 201 and a `Location` header pointing to the new resource

#### Scenario: Invalid hospital is rejected
- **WHEN** a client posts an invalid hospital creation payload
- **THEN** the server responds with HTTP 422 and a Problem Details JSON body

### Requirement: REST handlers delegate business behavior
REST handlers SHALL delegate Krankenhaus business behavior to a service layer instead of embedding persistence logic directly in handlers.

#### Scenario: Handler processes a create request
- **WHEN** a valid create request reaches the handler
- **THEN** the handler calls the service layer to create the hospital
