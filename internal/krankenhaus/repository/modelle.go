package repository

import "time"

type KrankenhausModel struct {
	ID                int                `gorm:"column:id;primaryKey;autoIncrement"`
	Version           int                `gorm:"column:version"`
	Name              string             `gorm:"column:name"`
	Mitarbeiteranzahl int                `gorm:"column:mitarbeiteranzahl"`
	Bettenanzahl      int                `gorm:"column:bettenanzahl"`
	Email             string             `gorm:"column:email"`
	Erzeugt           time.Time          `gorm:"column:erzeugt"`
	Aktualisiert      time.Time          `gorm:"column:aktualisiert"`
	Adresse           AdresseModel       `gorm:"foreignKey:KrankenhausID"`
	Fachbereiche      []FachbereichModel `gorm:"foreignKey:KrankenhausID"`
}

func (KrankenhausModel) TableName() string {
	return "krankenhaus"
}

type AdresseModel struct {
	ID            int    `gorm:"column:id;primaryKey;autoIncrement"`
	Strasse       string `gorm:"column:strasse"`
	Hausnummer    string `gorm:"column:hausnummer"`
	PLZ           string `gorm:"column:plz"`
	Ort           string `gorm:"column:ort"`
	KrankenhausID int    `gorm:"column:krankenhaus_id"`
}

func (AdresseModel) TableName() string {
	return "adresse"
}

type FachbereichModel struct {
	ID            int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `gorm:"column:name"`
	Beschreibung  string `gorm:"column:beschreibung"`
	Leitung       string `gorm:"column:leitung"`
	Anzahlaerzte  int    `gorm:"column:anzahlaerzte"`
	KrankenhausID int    `gorm:"column:krankenhaus_id"`
}

func (FachbereichModel) TableName() string {
	return "fachbereich"
}
