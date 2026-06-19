package service

import (
	"context"

	"github.com/johannesgt44/krankenhaus-ki/internal/krankenhaus/domain"
)

type Repository interface {
	SucheNachID(ctx context.Context, id int) (domain.Krankenhaus, error)
	Suche(ctx context.Context, suchparameter domain.Suchparameter) ([]domain.Krankenhaus, int64, error)
	Erstellen(ctx context.Context, krankenhaus domain.Krankenhaus) (int, error)
	Aktualisieren(ctx context.Context, id int, version int, krankenhaus domain.Krankenhaus) (int, error)
	Loeschen(ctx context.Context, id int) error
}

type KrankenhausService struct {
	repository Repository
}

func Neu(repository Repository) *KrankenhausService {
	return &KrankenhausService{repository: repository}
}

func (s *KrankenhausService) SucheNachID(ctx context.Context, id int) (domain.Krankenhaus, error) {
	return s.repository.SucheNachID(ctx, id)
}

func (s *KrankenhausService) Suche(ctx context.Context, suchparameter domain.Suchparameter) ([]domain.Krankenhaus, int64, error) {
	return s.repository.Suche(ctx, suchparameter)
}

func (s *KrankenhausService) Erstellen(ctx context.Context, krankenhaus domain.Krankenhaus) (int, error) {
	return s.repository.Erstellen(ctx, krankenhaus)
}

func (s *KrankenhausService) Aktualisieren(ctx context.Context, id int, version int, krankenhaus domain.Krankenhaus) (int, error) {
	return s.repository.Aktualisieren(ctx, id, version, krankenhaus)
}

func (s *KrankenhausService) Loeschen(ctx context.Context, id int) error {
	return s.repository.Loeschen(ctx, id)
}
