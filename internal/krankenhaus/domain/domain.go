package domain

import "time"

type Krankenhaus struct {
	ID                int
	Version           int
	Name              string
	Mitarbeiteranzahl int
	Bettenanzahl      int
	Email             string
	Erzeugt           time.Time
	Aktualisiert      time.Time
	Adresse           Adresse
	Fachbereiche      []Fachbereich
}

type Adresse struct {
	ID            int
	Strasse       string
	Hausnummer    string
	PLZ           string
	Ort           string
	KrankenhausID int
}

type Fachbereich struct {
	ID            int
	Name          string
	Beschreibung  string
	Leitung       string
	Anzahlaerzte  int
	KrankenhausID int
}

type Suchparameter struct {
	Name      string
	Page      int
	Size      int
	CountOnly bool
}
