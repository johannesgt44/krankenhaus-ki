package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	krankenhausrest "github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/rest"
	"github.com/johannesgt44/krankenhaus-ki/internal/problem"
)

type Option func(*optionen)

type optionen struct {
	schreibschutz func(http.Handler) http.Handler
}

func MitSchreibschutz(middleware func(http.Handler) http.Handler) Option {
	return func(optionen *optionen) {
		optionen.schreibschutz = middleware
	}
}

func Neu(dienst krankenhausrest.Dienst, opts ...Option) http.Handler {
	einstellungen := optionen{}
	for _, opt := range opts {
		opt(&einstellungen)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Mount("/rest/krankenhaus", krankenhausrest.NeuerRouter(
		dienst,
		krankenhausrest.MitSchreibschutz(einstellungen.schreibschutz),
	))
	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		problem.Schreiben(w, http.StatusNotFound, "Die angeforderte Ressource wurde nicht gefunden.")
	})

	return router
}
