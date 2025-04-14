package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error) {
	args := m.Called(ctx, productType, pvzID)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockService) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockService) GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.PVZWithReceptions), args.Error(1)
}

func (m *MockService) CreatePVZ(ctx context.Context, city string) (models.PVZ, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(models.PVZ), args.Error(1)
}
func (m *MockService) DummyLogin(ctx context.Context, req models.UserLoginReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockService) RegisterUser(ctx context.Context, req models.UserRegisterReq) (models.UserRegisterResp, error) {
	if req.Role == "admin" {
		return models.UserRegisterResp{}, errors.New("недопустимая роль")
	}

	if req.Email == "existing@example.com" {
		return models.UserRegisterResp{}, apperrors.ErrEmailAlreadyExists
	}

	return models.UserRegisterResp{
		ID:    uuid.New(),
		Email: req.Email,
		Role:  req.Role,
	}, nil
}

func (m *MockService) LoginUser(ctx context.Context, req models.UserLoginReq) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func (m *MockService) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func (m *MockService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}
