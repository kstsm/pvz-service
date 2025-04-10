package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
)

func (s Service) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	return s.repo.CreateReception(ctx, pvzID)
}

func (s Service) AddProductToReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error) {
	receptionID, err := s.repo.GetActiveReceptionByPVZ(ctx, pvzID)
	if err != nil {
		return models.Product{}, err
	}

	return s.repo.AddProduct(ctx, productType, receptionID)
}

func (s Service) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	return s.repo.DeleteLastProductInReception(ctx, pvzID)
}

func (s Service) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	reception, err := s.repo.CloseLastReception(ctx, pvzID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("ошибка при закрытии приёмки: %w", err)
	}

	return reception, nil
}
