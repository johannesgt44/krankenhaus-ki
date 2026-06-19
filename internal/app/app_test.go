package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/johannesgt44/krankenhaus-ki/internal/app"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/domain"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/service"
)

func TestHealth(t *testing.T) {
	server := httptest.NewServer(app.Neu(neuerFakeDienst()))
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("health request fehlgeschlagen: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status = %d, erwartet %d", resp.StatusCode, http.StatusOK)
	}
}

func TestKrankenhausCRUD(t *testing.T) {
	server := httptest.NewServer(app.Neu(neuerFakeDienst()))
	defer server.Close()

	payload := `{
		"name": "Staedtisches Klinikum Test",
		"mitarbeiteranzahl": 1200,
		"bettenanzahl": 450,
		"email": "test@klinikum.example",
		"adresse": {
			"strasse": "Teststrasse",
			"hausnummer": "1",
			"plz": "76133",
			"ort": "Karlsruhe"
		},
		"fachbereiche": [
			{"name": "Kardiologie", "beschreibung": "Herzmedizin", "leitung": "Dr. Test", "anzahlaerzte": 12}
		]
	}`

	createResp, err := http.Post(server.URL+"/rest/krankenhaus", "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("create request fehlgeschlagen: %v", err)
	}
	_ = createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create Status = %d, erwartet %d", createResp.StatusCode, http.StatusCreated)
	}
	location := createResp.Header.Get("Location")
	if location == "" {
		t.Fatal("Location Header fehlt")
	}

	readResp, err := http.Get(server.URL + location)
	if err != nil {
		t.Fatalf("read request fehlgeschlagen: %v", err)
	}
	defer func() {
		_ = readResp.Body.Close()
	}()
	if readResp.StatusCode != http.StatusOK {
		t.Fatalf("read Status = %d, erwartet %d", readResp.StatusCode, http.StatusOK)
	}
	if readResp.Header.Get("ETag") != `"0"` {
		t.Fatalf("ETag = %q, erwartet %q", readResp.Header.Get("ETag"), `"0"`)
	}

	updatePayload := `{
		"name": "Staedtisches Klinikum Aktualisiert",
		"mitarbeiteranzahl": 1300,
		"bettenanzahl": 470,
		"email": "aktualisiert@klinikum.example"
	}`
	updateReq, err := http.NewRequest(http.MethodPut, server.URL+location, bytes.NewBufferString(updatePayload))
	if err != nil {
		t.Fatalf("update request konnte nicht erstellt werden: %v", err)
	}
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("If-Match", `"0"`)
	updateResp, err := http.DefaultClient.Do(updateReq)
	if err != nil {
		t.Fatalf("update request fehlgeschlagen: %v", err)
	}
	_ = updateResp.Body.Close()
	if updateResp.StatusCode != http.StatusNoContent {
		t.Fatalf("update Status = %d, erwartet %d", updateResp.StatusCode, http.StatusNoContent)
	}
	if updateResp.Header.Get("ETag") != `"1"` {
		t.Fatalf("Update ETag = %q, erwartet %q", updateResp.Header.Get("ETag"), `"1"`)
	}

	deleteReq, err := http.NewRequest(http.MethodDelete, server.URL+location, nil)
	if err != nil {
		t.Fatalf("delete request konnte nicht erstellt werden: %v", err)
	}
	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatalf("delete request fehlgeschlagen: %v", err)
	}
	_ = deleteResp.Body.Close()
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete Status = %d, erwartet %d", deleteResp.StatusCode, http.StatusNoContent)
	}
}

func TestUngueltigerCreateRequestLiefertProblemDetails(t *testing.T) {
	server := httptest.NewServer(app.Neu(neuerFakeDienst()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/rest/krankenhaus", "application/json", strings.NewReader(`{"name": ""}`))
	if err != nil {
		t.Fatalf("invalid create request fehlgeschlagen: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("Status = %d, erwartet %d", resp.StatusCode, http.StatusUnprocessableEntity)
	}
	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "application/problem+json") {
		t.Fatalf("Content-Type = %q, erwartet Problem Details", contentType)
	}
}

type fakeDienst struct {
	mu       sync.Mutex
	naechste int
	daten    map[int]domain.Krankenhaus
}

func neuerFakeDienst() *fakeDienst {
	return &fakeDienst{
		naechste: 1000,
		daten:    map[int]domain.Krankenhaus{},
	}
}

func (f *fakeDienst) SucheNachID(_ context.Context, id int) (domain.Krankenhaus, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	krankenhaus, ok := f.daten[id]
	if !ok {
		return domain.Krankenhaus{}, service.ErrNichtGefunden
	}
	return krankenhaus, nil
}

func (f *fakeDienst) Suche(_ context.Context, _ domain.Suchparameter) ([]domain.Krankenhaus, int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	krankenhaeuser := make([]domain.Krankenhaus, 0, len(f.daten))
	for _, krankenhaus := range f.daten {
		krankenhaeuser = append(krankenhaeuser, krankenhaus)
	}
	return krankenhaeuser, int64(len(krankenhaeuser)), nil
}

func (f *fakeDienst) Erstellen(_ context.Context, krankenhaus domain.Krankenhaus) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.naechste
	f.naechste++
	krankenhaus.ID = id
	krankenhaus.Version = 0
	krankenhaus.Erzeugt = time.Now().UTC()
	krankenhaus.Aktualisiert = krankenhaus.Erzeugt
	f.daten[id] = krankenhaus
	return id, nil
}

func (f *fakeDienst) Aktualisieren(_ context.Context, id int, version int, krankenhaus domain.Krankenhaus) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	vorhanden, ok := f.daten[id]
	if !ok {
		return 0, service.ErrNichtGefunden
	}
	if vorhanden.Version != version {
		return 0, service.ErrVersionVeraltet
	}
	vorhanden.Version++
	vorhanden.Name = krankenhaus.Name
	vorhanden.Mitarbeiteranzahl = krankenhaus.Mitarbeiteranzahl
	vorhanden.Bettenanzahl = krankenhaus.Bettenanzahl
	vorhanden.Email = krankenhaus.Email
	vorhanden.Aktualisiert = time.Now().UTC()
	f.daten[id] = vorhanden
	return vorhanden.Version, nil
}

func (f *fakeDienst) Loeschen(_ context.Context, id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.daten, id)
	return nil
}

func TestListeUndCountOnly(t *testing.T) {
	server := httptest.NewServer(app.Neu(neuerFakeDienst()))
	defer server.Close()

	resp, err := http.Get(server.URL + "/rest/krankenhaus?count-only")
	if err != nil {
		t.Fatalf("count-only request fehlgeschlagen: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status = %d, erwartet %d", resp.StatusCode, http.StatusOK)
	}
	var body map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("JSON konnte nicht gelesen werden: %v", err)
	}
	if body["count"] != 0 {
		t.Fatalf("count = %d, erwartet 0", body["count"])
	}
}
