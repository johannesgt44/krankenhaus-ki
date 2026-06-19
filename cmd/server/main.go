package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/johannesgt44/krankenhaus-ki/internal/app"
	"github.com/johannesgt44/krankenhaus-ki/internal/config"
	"github.com/johannesgt44/krankenhaus-ki/internal/database"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/repository"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/service"
)

func main() {
	if err := starten(); err != nil {
		log.Fatalf("Server konnte nicht gestartet werden: %v", err)
	}
}

func starten() error {
	konfig := config.Laden()
	flag.StringVar(&konfig.ServerAdresse, "addr", konfig.ServerAdresse, "Adresse, auf der der Server lauscht")
	flag.BoolVar(&konfig.Datenbank.Init, "db-init", konfig.Datenbank.Init, "Datenbank beim Start neu erzeugen und befuellen")
	flag.Parse()

	db, err := database.Oeffnen(konfig.Datenbank)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer func() {
		_ = sqlDB.Close()
	}()

	repository := repository.Neu(db)
	dienst := service.Neu(repository)
	handler := app.Neu(dienst)

	server := &http.Server{
		Addr:              konfig.ServerAdresse,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Printf("Server lauscht auf %s", konfig.ServerAdresse)
		errChan <- server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		log.Printf("Server wird beendet: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	case err := <-errChan:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
