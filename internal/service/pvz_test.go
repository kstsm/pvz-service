package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreatePVZ(t *testing.T) {
	mockRepo := new(MockRepo)

	mockRepo.On("CreatePVZ", mock.Anything, "Moscow").Return(models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Moscow",
	}, nil)

	service := Service{repo: mockRepo}

	pvz, err := service.CreatePVZ(context.Background(), "Moscow")

	assert.NoError(t, err)
	assert.NotNil(t, pvz)
	assert.Equal(t, "Moscow", pvz.City)

	mockRepo.AssertExpectations(t)
}

func TestGetPVZList(t *testing.T) {
	mockRepo := new(MockRepo)

	mockRepo.On("GetPVZList", mock.Anything, models.PVZFilterParams{
		StartDate: nil,
		EndDate:   nil,
		Page:      1,
		Limit:     10,
	}).Return([]models.PVZWithReceptions{
		{
			PVZ: models.PVZ{
				ID:               uuid.New(),
				RegistrationDate: time.Now(),
				City:             "Moscow",
			},
			Receptions: []models.Reception{},
		},
	}, nil)

	service := Service{repo: mockRepo}

	params := models.PVZFilterParams{
		Page:  1,
		Limit: 10,
	}
	pvzList, err := service.GetPVZList(context.Background(), params)

	assert.NoError(t, err)
	assert.NotNil(t, pvzList)
	assert.Len(t, pvzList, 1)
	assert.Equal(t, "Moscow", pvzList[0].PVZ.City)

	mockRepo.AssertExpectations(t)
}
