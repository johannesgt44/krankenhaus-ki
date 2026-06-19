# Programmierworkshop am 19.6.2026

## Namen

Johannes Goette und Marc Ulm

## Link zum Git-Repository

https://github.com/johannesgt44/krankenhaus-ki

## Start

Der Server ist ein Go-REST-Server fuer die Krankenhaus-Domaene.

Zuerst die PostgreSQL-Datenbank per Docker Compose starten:

```powershell
docker compose -f extras/compose/compose.yml up -d db
```

Dann den Go-Server starten und die Tabellen fuer Entwicklung/Test neu erzeugen und befuellen:

```powershell
go run ./cmd/server -db-init
```

Ohne DB-Reset:

```powershell
go run ./cmd/server
```

Der Server lauscht standardmaessig auf `http://localhost:8080`. Die Defaults sind `localhost:5432`, Datenbank `postgres`, Benutzer `postgres`, Passwort `p`.

Datenbank stoppen:

```powershell
docker compose -f extras/compose/compose.yml down
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

Die Implementierung nutzt eine Go-typische Struktur mit `cmd/server`, `internal/app`, `internal/config`, `internal/database`, `internal/krankenhaus/domain`, `internal/krankenhaus/rest`, `internal/krankenhaus/service`, `internal/krankenhaus/repository`, `internal/problem` und `extras` fuer begleitende Dateien.

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

Die Schema-Erzeugung fuer Entwicklung/Test erfolgt ueber SQL-Scripts in `internal/database/sql`, nicht ueber GORM `AutoMigrate`.

### OIDC mit Keycloak

Keycloak/OIDC ist standardmaessig aktiv. Beim normalen Serverstart bleiben `GET /health`, `GET /rest/krankenhaus` und `GET /rest/krankenhaus/{id}` public. `POST`, `PUT` und `DELETE` unter `/rest/krankenhaus` benoetigen einen gueltigen Bearer-Token mit der Rolle `admin`.

Keycloak lokal starten:

```powershell
docker compose -f extras/compose/compose.yml up -d keycloak
```

Keycloak ist danach unter `http://localhost:8880` erreichbar.

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

## Beispiel Prompts/Requests an KI-Agent/en

- `openspec-explore: neuen Server aufsetzen`
- `Nutze openspec-rest-council.`
- `Sprache: go`
- `Framework: entscheide welches das bes passende framework ist fuer rest`
- `nutze zur implementierung eine go typische struktur bei rest schnittstellen und nutze deutsche syntax`
- `openspec-apply-change`
- `openspec-propose oidc mithilfe von Keycloak soll eingebaut werden. Bei allem außer beim Suchen Was brauchst du von mir alles damit du starten kannst`
