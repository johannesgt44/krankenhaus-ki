## ADDED Requirements

### Requirement: Integration test starts the REST app
The system SHALL include at least one integration test that starts the Go REST handler or server in a test context.

#### Scenario: Test app starts
- **WHEN** the integration test creates the app
- **THEN** the app can serve HTTP requests without a manual external server process

### Requirement: Integration test covers create flow
The system SHALL include an integration test for the hospital create REST flow.

#### Scenario: Create test succeeds
- **WHEN** the test posts a valid hospital payload
- **THEN** the response status is HTTP 201 and a `Location` header is present

### Requirement: Integration test covers read flow
The system SHALL include an integration test for reading a hospital through REST.

#### Scenario: Read test succeeds
- **WHEN** the test requests an existing hospital by ID
- **THEN** the response status is HTTP 200 and the response contains hospital JSON

### Requirement: Integration test covers invalid create input
The system SHALL include an integration test for invalid create input.

#### Scenario: Invalid create test succeeds
- **WHEN** the test posts an invalid hospital payload
- **THEN** the response status is HTTP 422 and the response content type is Problem Details JSON

### Requirement: Integration test covers update flow
The system SHALL include an integration test for updating a hospital through REST.

#### Scenario: Update test succeeds
- **WHEN** the test updates an existing hospital with a matching `If-Match` header
- **THEN** the response status is HTTP 204 and an `ETag` header is present

### Requirement: Integration test covers delete flow
The system SHALL include an integration test for deleting a hospital through REST.

#### Scenario: Delete test succeeds
- **WHEN** the test deletes an existing hospital
- **THEN** the response status is HTTP 204 and the hospital is no longer returned as an existing resource
