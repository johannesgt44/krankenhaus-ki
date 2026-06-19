## ADDED Requirements

### Requirement: Server connects to PostgreSQL through configuration
The system SHALL create PostgreSQL connections from documented configuration values.

#### Scenario: Database configuration is valid
- **WHEN** the server starts with valid PostgreSQL configuration
- **THEN** it can create a database connection pool

#### Scenario: Database configuration is invalid
- **WHEN** the server starts with invalid PostgreSQL configuration
- **THEN** startup fails with a clear error

### Requirement: Repository hides SQL access from handlers
The system SHALL access Krankenhaus persistence through a repository boundary instead of direct SQL in HTTP handlers.

#### Scenario: Read by ID uses repository
- **WHEN** the service reads a hospital by ID
- **THEN** it calls the repository to load the hospital data

#### Scenario: Create uses repository
- **WHEN** the service creates a hospital
- **THEN** it calls the repository to insert the hospital data and return the new ID

### Requirement: Persistence maps missing rows to domain errors
The system SHALL map missing database rows to a domain-level not-found error.

#### Scenario: Missing hospital is requested
- **WHEN** the repository does not find a hospital for an ID
- **THEN** the service exposes a not-found outcome that the REST layer maps to HTTP 404
