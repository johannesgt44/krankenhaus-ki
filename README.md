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

Die Schema-Erzeugung fuer Entwicklung/Test erfolgt ueber SQL-Scripts in `internal/database/sql`, nicht ueber GORM `AutoMigrate`.

### Optional: OIDC mit Keycloak

Keycloak/OIDC ist standardmaessig aktiv. Beim normalen Serverstart bleiben `GET /health`, `GET /rest/krankenhaus` und `GET /rest/krankenhaus/{id}` public. `POST`, `PUT` und `DELETE` unter `/rest/krankenhaus` benoetigen einen gueltigen Bearer-Token mit der Rolle `admin`.

Keycloak lokal starten:

```powershell
docker compose -f extras/compose/compose.yml up -d keycloak
```

Keycloak ist danach unter `http://localhost:8880` erreichbar. Die lokale Admin-Anmeldung lautet Benutzer `tmp`, Passwort `p`.

Lokale Keycloak-Konfiguration fuer die Entwicklung:

- Realm: `python`
- Client: `python-client`
- Client-Rolle: `admin`
- Testuser mit Admin-Rolle: z.B. `admin` mit Passwort `p`
- Testuser ohne Admin-Rolle: z.B. `user` mit Passwort `p`

Im Client `python-client` sollte fuer lokale curl-Tests Direct Access Grants aktiviert sein. Die Rolle `admin` wird als Client-Rolle am Client `python-client` angelegt und dem Admin-Testuser zugewiesen. Der Server liest Rollen aus `resource_access["python-client"].roles`; Realm-Rollen werden zusaetzlich als lokaler Fallback akzeptiert.

Server mit OIDC starten:

```powershell
$env:OIDC_ISSUER_URL="http://localhost:8880/realms/python"
$env:OIDC_CLIENT_ID="python-client"
$env:OIDC_REQUIRED_ROLE="admin"
go run ./cmd/server
```

Die drei Werte oben sind bereits Defaults. Fuer einen lokalen Start ohne Keycloak kann OIDC explizit abgeschaltet werden:

```powershell
$env:OIDC_ENABLED="false"
go run ./cmd/server
```

Token fuer den Admin-Testuser holen:

```powershell
$tokenResponse = Invoke-RestMethod -Method Post "http://localhost:8880/realms/python/protocol/openid-connect/token" -ContentType "application/x-www-form-urlencoded" -Body "grant_type=password&client_id=python-client&username=admin&password=p"
$token = $tokenResponse.access_token
```

Geschuetzten Endpoint mit Token aufrufen:

```powershell
Invoke-RestMethod -Method Post "http://localhost:8080/rest/krankenhaus" -Headers @{ Authorization = "Bearer $token" } -ContentType "application/json" -Body '{"name":"Staedtisches Klinikum Test","mitarbeiteranzahl":1200,"bettenanzahl":450,"email":"test@klinikum.example","adresse":{"strasse":"Teststrasse","hausnummer":"1","plz":"76133","ort":"Karlsruhe"},"fachbereiche":[{"name":"Kardiologie","beschreibung":"Herzmedizin","leitung":"Dr. Test","anzahlaerzte":12}]}'
```

Die passenden Bruno-Requests liegen unter `extras/bruno/krankenhaus`. Dort kann der Header `Authorization: Bearer <token>` fuer geschuetzte REST-Requests gesetzt werden.

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

- `openspec-explore: neuen Server aufsetzen`
- `Nutze openspec-rest-council.`
- `Sprache: go`
- `Framework: entscheide welches das bes passende framework ist fuer rest`
- `nutze zur implementierung eine go typische struktur bei rest schnittstellen und nutze deutsche syntax`
- `openspec-apply-change`
- `openspec-propose oidc mithilfe von Keycloak soll eingebaut werden. Bei allem außer beim Suchen Was brauchst du von mir alles damit du starten kannst`
