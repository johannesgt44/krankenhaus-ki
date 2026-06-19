package rest

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

const maxTextLaenge = 40

var plzRegex = regexp.MustCompile(`^\d{5}$`)

func validiereNeu(dto KrankenhausNeuDTO) []string {
	fehler := validiereBasis(dto.Name, dto.Mitarbeiteranzahl, dto.Bettenanzahl, dto.Email)
	fehler = append(fehler, validiereAdresse(dto.Adresse)...)
	for index, fachbereich := range dto.Fachbereiche {
		if strings.TrimSpace(fachbereich.Name) == "" {
			fehler = append(fehler, fmt.Sprintf("fachbereiche[%d].name ist erforderlich", index))
		}
		if len(fachbereich.Name) > maxTextLaenge {
			fehler = append(fehler, fmt.Sprintf("fachbereiche[%d].name darf hoechstens %d Zeichen lang sein", index, maxTextLaenge))
		}
		if fachbereich.Anzahlaerzte < 0 {
			fehler = append(fehler, fmt.Sprintf("fachbereiche[%d].anzahlaerzte darf nicht negativ sein", index))
		}
	}
	return fehler
}

func validiereUpdate(dto KrankenhausUpdateDTO) []string {
	return validiereBasis(dto.Name, dto.Mitarbeiteranzahl, dto.Bettenanzahl, dto.Email)
}

func validiereBasis(name string, mitarbeiteranzahl int, bettenanzahl int, email string) []string {
	var fehler []string
	if strings.TrimSpace(name) == "" {
		fehler = append(fehler, "name ist erforderlich")
	}
	if len(name) > maxTextLaenge {
		fehler = append(fehler, fmt.Sprintf("name darf hoechstens %d Zeichen lang sein", maxTextLaenge))
	}
	if mitarbeiteranzahl < 0 {
		fehler = append(fehler, "mitarbeiteranzahl darf nicht negativ sein")
	}
	if bettenanzahl < 0 {
		fehler = append(fehler, "bettenanzahl darf nicht negativ sein")
	}
	if strings.TrimSpace(email) == "" {
		fehler = append(fehler, "email ist erforderlich")
	} else if _, err := mail.ParseAddress(email); err != nil {
		fehler = append(fehler, "email ist ungueltig")
	}
	if len(email) > maxTextLaenge {
		fehler = append(fehler, fmt.Sprintf("email darf hoechstens %d Zeichen lang sein", maxTextLaenge))
	}
	return fehler
}

func validiereAdresse(adresse AdresseDTO) []string {
	var fehler []string
	if strings.TrimSpace(adresse.Strasse) == "" {
		fehler = append(fehler, "adresse.strasse ist erforderlich")
	}
	if strings.TrimSpace(adresse.Hausnummer) == "" {
		fehler = append(fehler, "adresse.hausnummer ist erforderlich")
	}
	if !plzRegex.MatchString(adresse.PLZ) {
		fehler = append(fehler, "adresse.plz muss aus genau fuenf Ziffern bestehen")
	}
	if strings.TrimSpace(adresse.Ort) == "" {
		fehler = append(fehler, "adresse.ort ist erforderlich")
	}
	return fehler
}
