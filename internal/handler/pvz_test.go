package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetListPVZ(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    func(*MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное получение списка ПВЗ",
			requestBody: map[string]int{
				"page":  1,
				"limit": 10,
			},
			mockService: func(m *MockService) {
				m.On("GetPVZList", mock.Anything, models.PVZFilterParams{
					Page:  1,
					Limit: 10,
				}).
					Return([]models.PVZWithReceptions{
						{
							PVZ: models.PVZ{
								ID:               uuid.New(),
								RegistrationDate: time.Now(),
								City:             "Москва",
							},
							Receptions: []models.Reception{
								{
									ID:       uuid.New(),
									DateTime: time.Now(),
									PVZID:    uuid.New(),
									Status:   "in_progress",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"city":"Москва"`,
		},
		{
			name: "Ошибка сервиса при получении списка ПВЗ",
			requestBody: map[string]int{
				"page":  1,
				"limit": 10,
			},
			mockService: func(m *MockService) {
				m.On("GetPVZList", mock.Anything, models.PVZFilterParams{
					Page:  1,
					Limit: 10,
				}).
					Return([]models.PVZWithReceptions(nil), errors.New("ошибка сервиса"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"Не удалось получить список ПВЗ"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			if tt.mockService != nil {
				tt.mockService(mockService)
			}

			handler := Handler{
				service: mockService,
			}

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.getListPVZ(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			body := buf.String()

			assert.Contains(t, body, tt.expectedBody)
		})
	}
}

func TestCreatePVZHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    func(*MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное создание ПВЗ",
			requestBody: map[string]string{
				"city": "Москва",
			},
			mockService: func(m *MockService) {
				m.On("CreatePVZ", mock.Anything, "Москва").
					Return(models.PVZ{
						ID:               uuid.New(),
						City:             "Москва",
						RegistrationDate: time.Now(),
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"city":"Москва"`,
		},
		{
			name: "Недопустимый город",
			requestBody: map[string]string{
				"city": "Тула",
			},
			mockService:    func(m *MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Данный город пока недоступен"`,
		},
		{
			name:           "Некорректный JSON",
			requestBody:    `{"city": Москва`,
			mockService:    func(m *MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Невалидный JSON"`,
		},
		{
			name: "Ошибка сервиса",
			requestBody: map[string]string{
				"city": "Казань",
			},
			mockService: func(m *MockService) {
				m.On("CreatePVZ", mock.Anything, "Казань").
					Return(models.PVZ{}, errors.New("ошибка сервиса"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"Не удалось создать ПВЗ"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			if tt.mockService != nil {
				tt.mockService(mockService)
			}

			handler := Handler{
				service: mockService,
			}

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.createPVZHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			body := buf.String()

			assert.Contains(t, body, tt.expectedBody)
		})
	}
}
