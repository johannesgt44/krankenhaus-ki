package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/domain"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/service"
	"github.com/johannesgt44/krankenhaus-ki/internal/problem"
)

type Dienst interface {
	SucheNachID(ctx context.Context, id int) (domain.Krankenhaus, error)
	Suche(ctx context.Context, suchparameter domain.Suchparameter) ([]domain.Krankenhaus, int64, error)
	Erstellen(ctx context.Context, krankenhaus domain.Krankenhaus) (int, error)
	Aktualisieren(ctx context.Context, id int, version int, krankenhaus domain.Krankenhaus) (int, error)
	Loeschen(ctx context.Context, id int) error
}

type Handler struct {
	dienst Dienst
}

type RouterOption func(*routerOptionen)

type routerOptionen struct {
	schreibschutz func(http.Handler) http.Handler
}

func MitSchreibschutz(middleware func(http.Handler) http.Handler) RouterOption {
	return func(optionen *routerOptionen) {
		optionen.schreibschutz = middleware
	}
}

func NeuerRouter(dienst Dienst, opts ...RouterOption) http.Handler {
	optionen := routerOptionen{}
	for _, opt := range opts {
		opt(&optionen)
	}

	handler := &Handler{dienst: dienst}
	router := chi.NewRouter()
	router.Get("/", handler.suche)
	router.Get("/{id}", handler.sucheNachID)
	if optionen.schreibschutz == nil {
		router.Post("/", handler.erstellen)
		router.Put("/{id}", handler.aktualisieren)
		router.Delete("/{id}", handler.loeschen)
	} else {
		router.With(optionen.schreibschutz).Post("/", handler.erstellen)
		router.With(optionen.schreibschutz).Put("/{id}", handler.aktualisieren)
		router.With(optionen.schreibschutz).Delete("/{id}", handler.loeschen)
	}
	return router
}

func (h *Handler) sucheNachID(w http.ResponseWriter, r *http.Request) {
	id, ok := idAusRequest(w, r)
	if !ok {
		return
	}

	krankenhaus, err := h.dienst.SucheNachID(r.Context(), id)
	if err != nil {
		h.schreibeFehler(w, err)
		return
	}

	etag := versionZuETag(krankenhaus.Version)
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)
	schreibeJSON(w, http.StatusOK, krankenhausZuDTO(krankenhaus))
}

func (h *Handler) suche(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := queryInt(query.Get("page"), 0)
	size := queryInt(query.Get("size"), 20)
	countOnly := false
	if _, ok := query["count-only"]; ok {
		countOnly = true
	}

	krankenhaeuser, anzahl, err := h.dienst.Suche(r.Context(), domain.Suchparameter{
		Name:      query.Get("name"),
		Page:      page,
		Size:      size,
		CountOnly: countOnly,
	})
	if err != nil {
		h.schreibeFehler(w, err)
		return
	}

	if countOnly {
		schreibeJSON(w, http.StatusOK, map[string]int64{"count": anzahl})
		return
	}

	content := make([]KrankenhausDTO, 0, len(krankenhaeuser))
	for _, krankenhaus := range krankenhaeuser {
		content = append(content, krankenhausZuDTO(krankenhaus))
	}
	schreibeJSON(w, http.StatusOK, SeiteDTO{
		Content:       content,
		TotalElements: anzahl,
		Page:          page,
		Size:          size,
	})
}

func (h *Handler) erstellen(w http.ResponseWriter, r *http.Request) {
	var dto KrankenhausNeuDTO
	if !leseJSON(w, r, &dto) {
		return
	}
	if fehler := validiereNeu(dto); len(fehler) > 0 {
		problem.Schreiben(w, http.StatusUnprocessableEntity, fehler)
		return
	}

	id, err := h.dienst.Erstellen(r.Context(), neuDTOZuDomain(dto))
	if err != nil {
		h.schreibeFehler(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%d", strings.TrimRight(r.URL.Path, "/"), id))
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) aktualisieren(w http.ResponseWriter, r *http.Request) {
	id, ok := idAusRequest(w, r)
	if !ok {
		return
	}
	version, ok := versionAusIfMatch(w, r)
	if !ok {
		return
	}

	var dto KrankenhausUpdateDTO
	if !leseJSON(w, r, &dto) {
		return
	}
	if fehler := validiereUpdate(dto); len(fehler) > 0 {
		problem.Schreiben(w, http.StatusUnprocessableEntity, fehler)
		return
	}

	neueVersion, err := h.dienst.Aktualisieren(r.Context(), id, version, updateDTOZuDomain(dto))
	if err != nil {
		h.schreibeFehler(w, err)
		return
	}
	w.Header().Set("ETag", versionZuETag(neueVersion))
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) loeschen(w http.ResponseWriter, r *http.Request) {
	id, ok := idAusRequest(w, r)
	if !ok {
		return
	}
	if err := h.dienst.Loeschen(r.Context(), id); err != nil {
		h.schreibeFehler(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) schreibeFehler(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNichtGefunden):
		problem.Schreiben(w, http.StatusNotFound, "Krankenhaus wurde nicht gefunden.")
	case errors.Is(err, service.ErrVersionVeraltet):
		problem.Schreiben(w, http.StatusPreconditionFailed, "Die angegebene Version ist veraltet.")
	case errors.Is(err, service.ErrEmailBereitsVorhanden):
		problem.Schreiben(w, http.StatusUnprocessableEntity, "Die Emailadresse ist bereits vergeben.")
	default:
		problem.Schreiben(w, http.StatusInternalServerError, "Interner Serverfehler.")
	}
}

func idAusRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		problem.Schreiben(w, http.StatusNotFound, "Die ID ist ungueltig.")
		return 0, false
	}
	return id, true
}

func versionAusIfMatch(w http.ResponseWriter, r *http.Request) (int, bool) {
	ifMatch := r.Header.Get("If-Match")
	if ifMatch == "" {
		problem.Schreiben(w, http.StatusPreconditionRequired, "Der Header If-Match ist erforderlich.")
		return 0, false
	}
	versionText := strings.Trim(ifMatch, "\"")
	version, err := strconv.Atoi(versionText)
	if err != nil || version < 0 {
		problem.Schreiben(w, http.StatusPreconditionFailed, "Der Header If-Match enthaelt keine gueltige Version.")
		return 0, false
	}
	return version, true
}

func leseJSON(w http.ResponseWriter, r *http.Request, ziel any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(ziel); err != nil {
		problem.Schreiben(w, http.StatusUnprocessableEntity, "Der JSON-Request ist ungueltig oder enthaelt unbekannte Felder.")
		return false
	}
	return true
}

func schreibeJSON(w http.ResponseWriter, statusCode int, daten any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(daten)
}

func versionZuETag(version int) string {
	return fmt.Sprintf("\"%d\"", version)
}

func queryInt(wert string, standard int) int {
	if wert == "" {
		return standard
	}
	zahl, err := strconv.Atoi(wert)
	if err != nil {
		return standard
	}
	return zahl
}
