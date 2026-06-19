package problem

import (
	"encoding/json"
	"net/http"
)

type Details struct {
	Title      string `json:"title"`
	StatusCode int    `json:"statusCode"`
	Detail     any    `json:"detail,omitempty"`
}

func Schreiben(w http.ResponseWriter, statusCode int, detail any) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(Details{
		Title:      titel(statusCode),
		StatusCode: statusCode,
		Detail:     detail,
	})
}

func titel(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Fehlerhafte Anfrage"
	case http.StatusUnauthorized:
		return "Nicht angemeldet"
	case http.StatusForbidden:
		return "Nicht erlaubt"
	case http.StatusNotFound:
		return "Nicht gefunden"
	case http.StatusPreconditionFailed:
		return "Vorbedingung fehlgeschlagen"
	case http.StatusUnprocessableEntity:
		return "Nicht verarbeitbarer Inhalt"
	case http.StatusPreconditionRequired:
		return "Vorbedingung erforderlich"
	default:
		return "Client-Fehler"
	}
}
