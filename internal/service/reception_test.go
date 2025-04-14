package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateReception(t *testing.T) {
	mockRepo := new(MockRepo)
	service := Service{repo: mockRepo}

	pvzID := uuid.New()
	expectedReception := models.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "created",
	}

	mockRepo.On("CreateReception", mock.Anything, pvzID).Return(expectedReception, nil)

	reception, err := service.CreateReception(context.Background(), pvzID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReception, reception)
	mockRepo.AssertExpectations(t)
}

func TestAddProductToActiveReception(t *testing.T) {
	mockRepo := new(MockRepo)
	service := Service{repo: mockRepo}

	pvzID := uuid.New()
	productType := "product_type"
	expectedProduct := models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        productType,
		ReceptionID: pvzID,
	}

	mockRepo.On("AddProductToActiveReception", mock.Anything, productType, pvzID).Return(expectedProduct, nil)

	product, err := service.AddProductToActiveReception(context.Background(), productType, pvzID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestDeleteLastProductInReception(t *testing.T) {
	mockRepo := new(MockRepo)
	service := Service{repo: mockRepo}

	pvzID := uuid.New()

	mockRepo.On("DeleteLastProductInReception", mock.Anything, pvzID).Return(nil)

	err := service.DeleteLastProductInReception(context.Background(), pvzID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCloseLastReception(t *testing.T) {
	mockRepo := new(MockRepo)
	service := Service{repo: mockRepo}

	pvzID := uuid.New()
	expectedReception := models.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "closed",
	}

	mockRepo.On("CloseLastReception", mock.Anything, pvzID).Return(expectedReception, nil)

	reception, err := service.CloseLastReception(context.Background(), pvzID)

	assert.NoError(t, err)
	assert.Equal(t, expectedReception, reception)
	mockRepo.AssertExpectations(t)
}

func TestCloseLastReceptionError(t *testing.T) {
	mockRepo := new(MockRepo)
	service := Service{repo: mockRepo}

	pvzID := uuid.New()

	mockRepo.On("CloseLastReception", mock.Anything, pvzID).Return(models.Reception{}, errors.New("db error"))

	reception, err := service.CloseLastReception(context.Background(), pvzID)

	assert.Error(t, err)
	assert.Equal(t, models.Reception{}, reception)
	mockRepo.AssertExpectations(t)
}
