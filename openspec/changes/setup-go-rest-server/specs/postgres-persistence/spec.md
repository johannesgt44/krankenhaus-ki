## ADDED Requirements

### Requirement: Server connects to PostgreSQL through configuration
The system SHALL create PostgreSQL connections from documented configuration values.

#### Scenario: Database configuration is valid
- **WHEN** the server starts with valid PostgreSQL configuration
- **THEN** it can create a database connection pool

#### Scenario: Database configuration is invalid
- **WHEN** the server starts with invalid PostgreSQL configuration
- **THEN** startup fails with a clear error

### Requirement: Persistence uses GORM mappings
The system SHALL use GORM models to map the existing PostgreSQL tables.

#### Scenario: Krankenhaus model is mapped
- **WHEN** the repository accesses Krankenhaus data
- **THEN** it uses a GORM model mapped to the `krankenhaus` table and its existing columns

#### Scenario: Related models are mapped
- **WHEN** the repository loads or saves addresses and departments
- **THEN** it uses GORM models mapped to `adresse` and `fachbereich`

### Requirement: Development startup can recreate and seed the database
The system SHALL provide a configurable development/test startup path that recreates the known schema and inserts seed data.

#### Scenario: Development initialization is enabled
- **WHEN** the server starts with database initialization enabled
- **THEN** it recreates the known tables and inserts seed data before accepting REST requests

#### Scenario: Development initialization is disabled
- **WHEN** the server starts with database initialization disabled
- **THEN** it connects to the existing database without dropping or recreating tables

#### Scenario: Schema reset uses SQL
- **WHEN** the database is recreated for development or tests
- **THEN** the schema is created from SQL that matches the supplied PostgreSQL schema

### Requirement: Repository hides SQL access from handlers
The system SHALL access Krankenhaus persistence through a repository boundary instead of direct SQL in HTTP handlers.

#### Scenario: Read by ID uses repository
- **WHEN** the service reads a hospital by ID
- **THEN** it calls the repository to load the hospital data

#### Scenario: Create uses repository
- **WHEN** the service creates a hospital
- **THEN** it calls the repository to insert the hospital data and return the new ID

#### Scenario: Update uses repository
- **WHEN** the service updates a hospital
- **THEN** it calls the repository to persist changes and increment or return the new version

#### Scenario: Delete uses repository
- **WHEN** the service deletes a hospital
- **THEN** it calls the repository to delete the hospital and dependent rows through the database relationship

### Requirement: Persistence maps missing rows to domain errors
The system SHALL map missing database rows to a domain-level not-found error.

#### Scenario: Missing hospital is requested
- **WHEN** the repository does not find a hospital for an ID
- **THEN** the service exposes a not-found outcome that the REST layer maps to HTTP 404
