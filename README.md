# Programmierworkshop am 19.6.2026

## Namen

Johannes Goette und Marc Ulm

## Link zum Git-Repository

https://github.com/johannesgt44/krankenhaus-ki

## Start

Der Server ist ein Go-REST-Server fuer die Krankenhaus-Domaene.

Zuerst die PostgreSQL-Datenbank per Docker Compose starten:

```powershell
docker compose up -d db
```

Dann den Go-Server starten und die Tabellen fuer Entwicklung/Test neu erzeugen und befuellen:

```powershell
go run ./cmd/server -db-init
```

Wenn `go` nicht im `PATH` ist:

```powershell
& "C:\Program Files\Go\bin\go.exe" run ./cmd/server -db-init
```

Ohne DB-Reset:

```powershell
go run ./cmd/server
```

Der Server lauscht standardmaessig auf `http://localhost:8080`. Die Defaults sind `localhost:5432`, Datenbank `postgres`, Benutzer `postgres`, Passwort `p`.

Datenbank stoppen:

```powershell
docker compose down
```

## KI-Werkzeuge

Spec-Driven Development Framework: OpenSpec (https://github.com/Fission-AI/OpenSpec/)

### Agenten

Codex / ChatGPT

## Frameworks und Bibliotheken

### REST-Schnittstelle

- Sprache: Go
- Router: `github.com/go-chi/chi/v5`
- Einstiegspunkt: `cmd/server`
- REST-Prefix: `/rest/krankenhaus`
- Health-Endpoint: `/health`

Unterstuetzte REST-Operationen:

- `GET /rest/krankenhaus`
- `GET /rest/krankenhaus/{id}`
- `POST /rest/krankenhaus`
- `PUT /rest/krankenhaus/{id}`
- `DELETE /rest/krankenhaus/{id}`

Die Implementierung nutzt eine Go-typische Struktur mit `internal/app`, `internal/konfiguration`, `internal/krankenhaus/rest`, `internal/krankenhaus/service`, `internal/krankenhaus/repository`, `internal/datenbank` und `internal/problem`.

### Deutsche Domaenensyntax

Domain-Begriffe und JSON-Felder sind deutsch bzw. deutsch mit ASCII-Schreibweise:

- `Krankenhaus`
- `Adresse`
- `Fachbereich`
- `mitarbeiteranzahl`
- `bettenanzahl`
- `strasse`
- `hausnummer`
- `plz`
- `fachbereiche`

Fehlerantworten werden als Problem Details mit deutschen Details ausgeliefert.

### Validierung

Create- und Update-Payloads werden vor der Persistenz validiert. Unbekannte JSON-Felder werden abgelehnt. Fehlerhafte Requests liefern HTTP `422` mit `application/problem+json`.

Validiert werden u.a.:

- Pflichtfelder
- Emailformat
- nicht-negative Zahlen fuer Betten/Mitarbeiter/Aerzte
- fuenfstellige PLZ
- einfache Textlaengen

### OR-Mapping fuer PostgreSQL

ORM: GORM

- `gorm.io/gorm`
- `gorm.io/driver/postgres`

Die Modelle sind explizit auf die vorhandenen Tabellen gemappt:

- `krankenhaus`
- `adresse`
- `fachbereich`

Die Schema-Erzeugung fuer Entwicklung/Test erfolgt ueber SQL-Scripts in `internal/datenbank/sql`, nicht ueber GORM `AutoMigrate`.

### Optional: OIDC mit Keycloak

Keycloak/OIDC ist optional und in dieser ersten Umsetzung nicht erforderlich. Die REST-Struktur ist so aufgebaut, dass spaeter Middleware fuer Authentifizierung/Autorisierung ergaenzt werden kann.

### Einfacher Integrationstest

Die Integrationstests nutzen `net/http/httptest` und pruefen die REST-App ohne manuell gestarteten externen Server.

```powershell
go test ./...
```

Getestet werden:

- Health-Endpoint
- Create mit `201 Created` und `Location`
- Read mit `ETag`
- Update mit `If-Match` und neuem `ETag`
- Delete mit `204 No Content`
- Invalid Create mit `422` und Problem Details
- `count-only`

## Prompts/Requests an KI-Agent/en

- `/openspec-explore: neuen Server aufsetzen`
- `Nutze openspec-rest-council.`
- `Sprache: go`
- `Framework: entscheide welches das bes passende framework ist fuer rest`
- `Nutze danach OpenSpec, um proposal/design/tasks zu erstellen.`
- `Noch nicht implementieren, bis ich "apply" sage.`
- `nutze zur implementierung eine go typische struktur bei rest schnittstellen und nutze deutsche syntax`
- `ORM bitte GORM`
- `auch update und delete`
- `alles deutsch`
- `/openspec-apply-change`
