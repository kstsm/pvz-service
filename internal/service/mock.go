package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockRepo) AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error) {
	args := m.Called(ctx, productType, pvzID)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockRepo) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}

func (m *MockRepo) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockRepo) CreatePVZ(ctx context.Context, city string) (models.PVZ, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(models.PVZ), args.Error(1)
}
func (m *MockRepo) GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]models.PVZWithReceptions), args.Error(1)
}

func (m *MockRepo) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepo) CreateUser(ctx context.Context, user models.UserRegisterReq) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRepo) GetRoleByEmail(ctx context.Context, email string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}
