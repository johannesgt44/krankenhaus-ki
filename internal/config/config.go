package config

import (
	"fmt"
	"os"
	"strconv"
)

type Konfiguration struct {
	ServerAdresse string
	Datenbank     DatenbankKonfiguration
	OIDC          OIDCKonfiguration
}

type DatenbankKonfiguration struct {
	Host     string
	Port     int
	Name     string
	Benutzer string
	Passwort string
	SSLModus string
	Init     bool
}

type OIDCKonfiguration struct {
	Aktiv              bool
	IssuerURL          string
	ClientID           string
	ErforderlicheRolle string
}

func Laden() Konfiguration {
	return Konfiguration{
		ServerAdresse: envString("SERVER_ADRESSE", ":8080"),
		Datenbank: DatenbankKonfiguration{
			Host:     envString("DB_HOST", "localhost"),
			Port:     envInt("DB_PORT", 5432),
			Name:     envString("DB_NAME", "postgres"),
			Benutzer: envString("DB_USER", "postgres"),
			Passwort: envString("DB_PASSWORD", "p"),
			SSLModus: envString("DB_SSLMODE", "disable"),
			Init:     envBool("DB_INIT", false),
		},
		OIDC: OIDCKonfiguration{
			Aktiv:              envBool("OIDC_ENABLED", true),
			IssuerURL:          envString("OIDC_ISSUER_URL", "http://localhost:8880/realms/javascript"),
			ClientID:           envString("OIDC_CLIENT_ID", "javascript-client"),
			ErforderlicheRolle: envString("OIDC_REQUIRED_ROLE", "admin"),
		},
	}
}

func (k DatenbankKonfiguration) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		k.Host,
		k.Port,
		k.Benutzer,
		k.Passwort,
		k.Name,
		k.SSLModus,
	)
}

func envString(name string, standard string) string {
	wert := os.Getenv(name)
	if wert == "" {
		return standard
	}
	return wert
}

func envInt(name string, standard int) int {
	wert := os.Getenv(name)
	if wert == "" {
		return standard
	}
	zahl, err := strconv.Atoi(wert)
	if err != nil {
		return standard
	}
	return zahl
}

func envBool(name string, standard bool) bool {
	wert := os.Getenv(name)
	if wert == "" {
		return standard
	}
	parsed, err := strconv.ParseBool(wert)
	if err != nil {
		return standard
	}
	return parsed
}
