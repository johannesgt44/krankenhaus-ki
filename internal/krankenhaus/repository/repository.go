package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/domain"
	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/service"
	"gorm.io/gorm"
)

type KrankenhausRepository struct {
	db *gorm.DB
}

func Neu(db *gorm.DB) *KrankenhausRepository {
	return &KrankenhausRepository{db: db}
}

func (r *KrankenhausRepository) SucheNachID(ctx context.Context, id int) (domain.Krankenhaus, error) {
	var model KrankenhausModel
	err := r.db.WithContext(ctx).
		Preload("Adresse").
		Preload("Fachbereiche").
		First(&model, "id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Krankenhaus{}, service.ErrNichtGefunden
	}
	if err != nil {
		return domain.Krankenhaus{}, err
	}
	return modelZuDomain(model), nil
}

func (r *KrankenhausRepository) Suche(ctx context.Context, suchparameter domain.Suchparameter) ([]domain.Krankenhaus, int64, error) {
	query := r.db.WithContext(ctx).Model(&KrankenhausModel{})
	if suchparameter.Name != "" {
		query = query.Where("lower(name) LIKE ?", "%"+strings.ToLower(suchparameter.Name)+"%")
	}

	var anzahl int64
	if err := query.Count(&anzahl).Error; err != nil {
		return nil, 0, err
	}

	page := suchparameter.Page
	if page < 0 {
		page = 0
	}
	size := suchparameter.Size
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	var modelle []KrankenhausModel
	err := query.
		Preload("Adresse").
		Preload("Fachbereiche").
		Order("id ASC").
		Limit(size).
		Offset(page * size).
		Find(&modelle).
		Error
	if err != nil {
		return nil, 0, err
	}

	krankenhaeuser := make([]domain.Krankenhaus, 0, len(modelle))
	for _, model := range modelle {
		krankenhaeuser = append(krankenhaeuser, modelZuDomain(model))
	}
	return krankenhaeuser, anzahl, nil
}

func (r *KrankenhausRepository) Erstellen(ctx context.Context, krankenhaus domain.Krankenhaus) (int, error) {
	jetzt := time.Now().UTC()
	model := domainZuModel(krankenhaus)
	model.ID = 0
	model.Version = 0
	model.Erzeugt = jetzt
	model.Aktualisiert = jetzt

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			if istUniqueFehler(err) {
				return service.ErrEmailBereitsVorhanden
			}
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return model.ID, nil
}

func (r *KrankenhausRepository) Aktualisieren(ctx context.Context, id int, version int, krankenhaus domain.Krankenhaus) (int, error) {
	var neueVersion int
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var vorhanden KrankenhausModel
		if err := tx.First(&vorhanden, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return service.ErrNichtGefunden
			}
			return err
		}
		if vorhanden.Version != version {
			return service.ErrVersionVeraltet
		}

		neueVersion = vorhanden.Version + 1
		updates := map[string]any{
			"version":           neueVersion,
			"name":              krankenhaus.Name,
			"mitarbeiteranzahl": krankenhaus.Mitarbeiteranzahl,
			"bettenanzahl":      krankenhaus.Bettenanzahl,
			"email":             krankenhaus.Email,
			"aktualisiert":      time.Now().UTC(),
		}
		err := tx.Model(&KrankenhausModel{}).Where("id = ?", id).Updates(updates).Error
		if err != nil {
			if istUniqueFehler(err) {
				return service.ErrEmailBereitsVorhanden
			}
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return neueVersion, nil
}

func (r *KrankenhausRepository) Loeschen(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&KrankenhausModel{}).Error
}

func domainZuModel(krankenhaus domain.Krankenhaus) KrankenhausModel {
	fachbereiche := make([]FachbereichModel, 0, len(krankenhaus.Fachbereiche))
	for _, fachbereich := range krankenhaus.Fachbereiche {
		fachbereiche = append(fachbereiche, FachbereichModel{
			ID:           fachbereich.ID,
			Name:         fachbereich.Name,
			Beschreibung: fachbereich.Beschreibung,
			Leitung:      fachbereich.Leitung,
			Anzahlaerzte: fachbereich.Anzahlaerzte,
		})
	}

	return KrankenhausModel{
		ID:                krankenhaus.ID,
		Version:           krankenhaus.Version,
		Name:              krankenhaus.Name,
		Mitarbeiteranzahl: krankenhaus.Mitarbeiteranzahl,
		Bettenanzahl:      krankenhaus.Bettenanzahl,
		Email:             krankenhaus.Email,
		Erzeugt:           krankenhaus.Erzeugt,
		Aktualisiert:      krankenhaus.Aktualisiert,
		Adresse: AdresseModel{
			ID:         krankenhaus.Adresse.ID,
			Strasse:    krankenhaus.Adresse.Strasse,
			Hausnummer: krankenhaus.Adresse.Hausnummer,
			PLZ:        krankenhaus.Adresse.PLZ,
			Ort:        krankenhaus.Adresse.Ort,
		},
		Fachbereiche: fachbereiche,
	}
}

func modelZuDomain(model KrankenhausModel) domain.Krankenhaus {
	fachbereiche := make([]domain.Fachbereich, 0, len(model.Fachbereiche))
	for _, fachbereich := range model.Fachbereiche {
		fachbereiche = append(fachbereiche, domain.Fachbereich{
			ID:            fachbereich.ID,
			Name:          fachbereich.Name,
			Beschreibung:  fachbereich.Beschreibung,
			Leitung:       fachbereich.Leitung,
			Anzahlaerzte:  fachbereich.Anzahlaerzte,
			KrankenhausID: fachbereich.KrankenhausID,
		})
	}

	return domain.Krankenhaus{
		ID:                model.ID,
		Version:           model.Version,
		Name:              model.Name,
		Mitarbeiteranzahl: model.Mitarbeiteranzahl,
		Bettenanzahl:      model.Bettenanzahl,
		Email:             model.Email,
		Erzeugt:           model.Erzeugt,
		Aktualisiert:      model.Aktualisiert,
		Adresse: domain.Adresse{
			ID:            model.Adresse.ID,
			Strasse:       model.Adresse.Strasse,
			Hausnummer:    model.Adresse.Hausnummer,
			PLZ:           model.Adresse.PLZ,
			Ort:           model.Adresse.Ort,
			KrankenhausID: model.Adresse.KrankenhausID,
		},
		Fachbereiche: fachbereiche,
	}
}

func istUniqueFehler(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
