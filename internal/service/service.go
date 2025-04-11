package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/internal/repository"
	"github.com/kstsm/pvz-service/models"
)

type ServiceI interface {
	CreatePVZ(ctx context.Context, city string) (models.PVZ, error)
	CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error)
	AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error)
	DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error)
	GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error)
	RegisterUser(ctx context.Context, req models.UserRegisterReq) (models.UserRegisterResp, error)
	LoginUser(ctx context.Context, req models.UserLoginReq) (string, error)
}

type Service struct {
	repo repository.RepositoryI
}

func NewService(repo repository.RepositoryI) *Service {
	return &Service{
		repo: repo,
	}
}
