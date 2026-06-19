## ADDED Requirements

### Requirement: Create payload is validated before persistence
The system SHALL validate hospital creation JSON before calling persistence.

#### Scenario: Required fields are present
- **WHEN** a create request includes all required hospital fields with valid values
- **THEN** validation succeeds and the request can continue to the service layer

#### Scenario: Required field is missing
- **WHEN** a create request omits a required hospital field
- **THEN** validation fails with HTTP 422

### Requirement: Update payload is validated before persistence
The system SHALL validate hospital update JSON before calling persistence.

#### Scenario: Update fields are valid
- **WHEN** an update request includes valid updatable hospital fields
- **THEN** validation succeeds and the request can continue to the service layer

#### Scenario: Update field is invalid
- **WHEN** an update request includes an invalid updatable hospital field
- **THEN** validation fails with HTTP 422

### Requirement: Hospital scalar fields follow domain constraints
The system SHALL validate scalar hospital fields for type, length, and format before create or update.

#### Scenario: Email is invalid
- **WHEN** a create or update request contains an invalid email address
- **THEN** validation fails with HTTP 422

#### Scenario: Numeric count is negative
- **WHEN** a create or update request contains a negative bed count or employee count
- **THEN** validation fails with HTTP 422

### Requirement: Address fields are validated
The system SHALL validate address fields in a hospital create payload.

#### Scenario: Postal code has invalid format
- **WHEN** a create request contains a postal code that is not a five-digit value
- **THEN** validation fails with HTTP 422

### Requirement: Unknown JSON handling is deterministic
The system SHALL define whether unknown JSON fields are rejected or ignored and apply that behavior consistently.

#### Scenario: Unknown field is submitted
- **WHEN** a create or update request contains an unknown JSON field
- **THEN** the server applies the documented unknown-field behavior consistently
