package datenbank

import (
	"embed"
	"fmt"
	"strings"

	"github.com/johannesgt44/krankenhaus-ki/internal/konfiguration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed sql/schema.sql sql/seed.sql
var scripts embed.FS

func Oeffnen(konfig konfiguration.DatenbankKonfiguration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(konfig.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("datenbankverbindung fehlgeschlagen: %w", err)
	}

	if konfig.Init {
		if err := Initialisieren(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func Initialisieren(db *gorm.DB) error {
	for _, datei := range []string{"sql/schema.sql", "sql/seed.sql"} {
		inhalt, err := scripts.ReadFile(datei)
		if err != nil {
			return fmt.Errorf("sql-script %s konnte nicht gelesen werden: %w", datei, err)
		}
		if err := ausfuehren(db, string(inhalt)); err != nil {
			return fmt.Errorf("sql-script %s fehlgeschlagen: %w", datei, err)
		}
	}
	return nil
}

func ausfuehren(db *gorm.DB, script string) error {
	for _, statement := range strings.Split(script, ";") {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}
