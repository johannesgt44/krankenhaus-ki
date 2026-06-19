package rest

import (
	"strings"
	"time"

	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/domain"
)

type KrankenhausDTO struct {
	ID                int              `json:"id"`
	Name              string           `json:"name"`
	Mitarbeiteranzahl int              `json:"mitarbeiteranzahl"`
	Bettenanzahl      int              `json:"bettenanzahl"`
	Email             string           `json:"email"`
	Adresse           AdresseDTO       `json:"adresse"`
	Fachbereiche      []FachbereichDTO `json:"fachbereiche"`
	Erzeugt           time.Time        `json:"erzeugt"`
	Aktualisiert      time.Time        `json:"aktualisiert"`
}

type KrankenhausNeuDTO struct {
	Name              string           `json:"name"`
	Mitarbeiteranzahl int              `json:"mitarbeiteranzahl"`
	Bettenanzahl      int              `json:"bettenanzahl"`
	Email             string           `json:"email"`
	Adresse           AdresseDTO       `json:"adresse"`
	Fachbereiche      []FachbereichDTO `json:"fachbereiche"`
}

type KrankenhausUpdateDTO struct {
	Name              string `json:"name"`
	Mitarbeiteranzahl int    `json:"mitarbeiteranzahl"`
	Bettenanzahl      int    `json:"bettenanzahl"`
	Email             string `json:"email"`
}

type AdresseDTO struct {
	Strasse    string `json:"strasse"`
	Hausnummer string `json:"hausnummer"`
	PLZ        string `json:"plz"`
	Ort        string `json:"ort"`
}

type FachbereichDTO struct {
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung,omitempty"`
	Leitung      string `json:"leitung,omitempty"`
	Anzahlaerzte int    `json:"anzahlaerzte,omitempty"`
}

type SeiteDTO struct {
	Content       []KrankenhausDTO `json:"content"`
	TotalElements int64            `json:"totalElements"`
	Page          int              `json:"page"`
	Size          int              `json:"size"`
}

func krankenhausZuDTO(krankenhaus domain.Krankenhaus) KrankenhausDTO {
	fachbereiche := make([]FachbereichDTO, 0, len(krankenhaus.Fachbereiche))
	for _, fachbereich := range krankenhaus.Fachbereiche {
		fachbereiche = append(fachbereiche, FachbereichDTO{
			Name:         fachbereich.Name,
			Beschreibung: fachbereich.Beschreibung,
			Leitung:      fachbereich.Leitung,
			Anzahlaerzte: fachbereich.Anzahlaerzte,
		})
	}

	return KrankenhausDTO{
		ID:                krankenhaus.ID,
		Name:              krankenhaus.Name,
		Mitarbeiteranzahl: krankenhaus.Mitarbeiteranzahl,
		Bettenanzahl:      krankenhaus.Bettenanzahl,
		Email:             krankenhaus.Email,
		Erzeugt:           krankenhaus.Erzeugt,
		Aktualisiert:      krankenhaus.Aktualisiert,
		Adresse: AdresseDTO{
			Strasse:    krankenhaus.Adresse.Strasse,
			Hausnummer: krankenhaus.Adresse.Hausnummer,
			PLZ:        krankenhaus.Adresse.PLZ,
			Ort:        krankenhaus.Adresse.Ort,
		},
		Fachbereiche: fachbereiche,
	}
}

func neuDTOZuDomain(dto KrankenhausNeuDTO) domain.Krankenhaus {
	fachbereiche := make([]domain.Fachbereich, 0, len(dto.Fachbereiche))
	for _, fachbereich := range dto.Fachbereiche {
		fachbereiche = append(fachbereiche, domain.Fachbereich{
			Name:         strings.TrimSpace(fachbereich.Name),
			Beschreibung: strings.TrimSpace(fachbereich.Beschreibung),
			Leitung:      strings.TrimSpace(fachbereich.Leitung),
			Anzahlaerzte: fachbereich.Anzahlaerzte,
		})
	}

	return domain.Krankenhaus{
		Name:              strings.TrimSpace(dto.Name),
		Mitarbeiteranzahl: dto.Mitarbeiteranzahl,
		Bettenanzahl:      dto.Bettenanzahl,
		Email:             strings.TrimSpace(dto.Email),
		Adresse: domain.Adresse{
			Strasse:    strings.TrimSpace(dto.Adresse.Strasse),
			Hausnummer: strings.TrimSpace(dto.Adresse.Hausnummer),
			PLZ:        strings.TrimSpace(dto.Adresse.PLZ),
			Ort:        strings.TrimSpace(dto.Adresse.Ort),
		},
		Fachbereiche: fachbereiche,
	}
}

func updateDTOZuDomain(dto KrankenhausUpdateDTO) domain.Krankenhaus {
	return domain.Krankenhaus{
		Name:              strings.TrimSpace(dto.Name),
		Mitarbeiteranzahl: dto.Mitarbeiteranzahl,
		Bettenanzahl:      dto.Bettenanzahl,
		Email:             strings.TrimSpace(dto.Email),
	}
}
