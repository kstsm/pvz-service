package service

import (
	"context"
	"github.com/kstsm/pvz-service/models"
)

func (s Service) CreatePVZ(ctx context.Context, city string) (models.PVZ, error) {
	return s.repo.CreatePVZ(ctx, city)
}

func (s Service) GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error) {
	return s.repo.GetPVZList(ctx, params)
}
