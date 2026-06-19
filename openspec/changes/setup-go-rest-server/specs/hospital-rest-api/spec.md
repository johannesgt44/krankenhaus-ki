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

### Requirement: REST JSON uses German Krankenhaus field names
The system SHALL expose Krankenhaus JSON using German domain field names with ASCII spelling.

#### Scenario: Hospital JSON is returned
- **WHEN** a client reads a hospital
- **THEN** the JSON fields include German names such as `mitarbeiteranzahl`, `bettenanzahl`, `adresse`, `strasse`, `plz`, `ort`, and `fachbereiche`

#### Scenario: Hospital JSON is created
- **WHEN** a client posts a hospital creation payload
- **THEN** the server accepts the documented German JSON field names

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

### Requirement: Client can update a hospital
The system SHALL expose `PUT /rest/krankenhaus/{id}` to update an existing hospital from a valid JSON request body.

#### Scenario: Valid hospital is updated
- **WHEN** a client sends a valid update payload with a matching `If-Match` version
- **THEN** the server responds with HTTP 204 and a new `ETag` header

#### Scenario: Update is missing If-Match
- **WHEN** a client sends an update request without `If-Match`
- **THEN** the server responds with HTTP 428 and a Problem Details JSON body

#### Scenario: Update version is stale
- **WHEN** a client sends an update request with an outdated `If-Match` version
- **THEN** the server responds with HTTP 412 and a Problem Details JSON body

#### Scenario: Invalid update is rejected
- **WHEN** a client sends an invalid update payload
- **THEN** the server responds with HTTP 422 and a Problem Details JSON body

### Requirement: Client can delete a hospital
The system SHALL expose `DELETE /rest/krankenhaus/{id}` to delete an existing hospital.

#### Scenario: Existing hospital is deleted
- **WHEN** a client deletes an existing hospital ID
- **THEN** the server responds with HTTP 204 and no response body

#### Scenario: Missing hospital delete is idempotent
- **WHEN** a client deletes a missing hospital ID
- **THEN** the server responds with HTTP 204 and no response body

### Requirement: REST handlers delegate business behavior
REST handlers SHALL delegate Krankenhaus business behavior to a service layer instead of embedding persistence logic directly in handlers.

#### Scenario: Handler processes a create request
- **WHEN** a valid create request reaches the handler
- **THEN** the handler calls the service layer to create the hospital

#### Scenario: Handler processes an update request
- **WHEN** a valid update request reaches the handler
- **THEN** the handler calls the service layer to update the hospital

#### Scenario: Handler processes a delete request
- **WHEN** a delete request reaches the handler
- **THEN** the handler calls the service layer to delete the hospital
