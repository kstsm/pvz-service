package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/pvz-service/models"
)

type RepositoryI interface {
	CreatePVZ(ctx context.Context, city string) (models.PVZ, error)
	CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error)
	AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error)
	DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error)
	GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user models.UserRegisterReq) (uuid.UUID, error)
	GetRoleByEmail(ctx context.Context, email string) (string, string, error)
}

type Repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) RepositoryI {
	return &Repository{
		conn: conn,
	}
}
