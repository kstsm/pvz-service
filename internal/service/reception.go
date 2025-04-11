package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
)

func (s Service) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	return s.repo.CreateReception(ctx, pvzID)
}

func (s Service) AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error) {
	return s.repo.AddProductToActiveReception(ctx, productType, pvzID)
}

func (s Service) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	return s.repo.DeleteLastProductInReception(ctx, pvzID)
}

func (s Service) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	reception, err := s.repo.CloseLastReception(ctx, pvzID)
	if err != nil {
		return models.Reception{}, err
	}

	return reception, nil
}
