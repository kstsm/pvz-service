package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCloseLastReceptionHandler(t *testing.T) {
	pvzID := uuid.New()
	reception := models.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   "in_progress",
	}

	tests := []struct {
		name               string
		pvzIDParam         string
		userRole           string
		expectedStatusCode int
		expectedResponse   models.Error
		mockServiceReturn  models.Reception
		mockServiceError   error
	}{
		{
			name:               "Employee trying to close a reception that is already closed",
			pvzIDParam:         pvzID.String(),
			userRole:           "employee",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.Error{Message: "Приемка уже закрыта или не найдена"},
			mockServiceReturn:  reception,
			mockServiceError:   apperrors.ErrReceptionAlreadyClosed,
		},
		{
			name:               "Employee successfully closes reception",
			pvzIDParam:         pvzID.String(),
			userRole:           "employee",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   models.Error{},
			mockServiceReturn:  reception,
			mockServiceError:   nil,
		},
		{
			name:               "Reception not found for employee",
			pvzIDParam:         pvzID.String(),
			userRole:           "employee",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.Error{Message: "Приемка уже закрыта или не найдена"},
			mockServiceReturn:  reception,
			mockServiceError:   apperrors.ErrReceptionAlreadyClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockService.On("CloseLastReception", mock.Anything, pvzID).Return(tt.mockServiceReturn, tt.mockServiceError)

			h := Handler{
				service: mockService,
			}

			req := httptest.NewRequest(http.MethodPost, "/pvz/"+tt.pvzIDParam+"/close_last_reception", nil)
			req.Header.Set("Authorization", "Bearer valid_token")
			req = req.WithContext(context.WithValue(req.Context(), "userRole", tt.userRole))

			w := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/pvz/{pvzId}/close_last_reception", h.closeLastReceptionHandler)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)

			var errResp models.Error
			err := json.NewDecoder(w.Body).Decode(&errResp)
			if err != nil {
				t.Fatalf("Error decoding response: %v", err)
			}

			assert.Equal(t, tt.expectedResponse, errResp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteLastProductHandler(t *testing.T) {
	tests := []struct {
		name               string
		pvzID              string
		role               string
		expectedStatusCode int
		mockServiceError   error
		expectedResponse   models.Error
	}{
		{
			name:               "Authorized user, active reception, product exists",
			pvzID:              "8eabafc5-365b-4d43-8ac4-e7b70c8a4db2",
			role:               "employee",
			expectedStatusCode: http.StatusOK,
			mockServiceError:   nil,
		},
		{
			name:               "No active reception",
			pvzID:              "8eabafc5-365b-4d43-8ac4-e7b70c8a4db2",
			role:               "employee",
			expectedStatusCode: http.StatusBadRequest,
			mockServiceError:   apperrors.ErrNoActiveReception,
			expectedResponse:   models.Error{Message: "Нет активной приёмки для удаления товара"},
		},
		{
			name:               "No products to delete",
			pvzID:              "8eabafc5-365b-4d43-8ac4-e7b70c8a4db2",
			role:               "employee",
			expectedStatusCode: http.StatusBadRequest,
			mockServiceError:   apperrors.ErrNoProductToDelete,
			expectedResponse:   models.Error{Message: "Нет товаров для удаления в активной приёмке"},
		}, {
			name:               "Invalid UUID format",
			pvzID:              "invalid-uuid",
			role:               "employee",
			expectedStatusCode: http.StatusBadRequest,
			mockServiceError:   nil,
			expectedResponse:   models.Error{Message: "Неверный формат идентификатора ПВЗ"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := new(MockService)
			if tt.mockServiceError != nil {
				serviceMock.On("DeleteLastProductInReception", mock.Anything, mock.Anything).Return(tt.mockServiceError)
			} else {
				serviceMock.On("DeleteLastProductInReception", mock.Anything, mock.Anything).Return(nil)
			}

			handler := Handler{
				service: serviceMock,
			}

			req := httptest.NewRequest(http.MethodPost, "/pvz/"+tt.pvzID+"/delete_last_product", nil)
			req.Header.Set("Authorization", "Bearer "+generateMockToken(tt.role))

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("pvzId", tt.pvzID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			rr := httptest.NewRecorder()

			handler.deleteLastProductHandler(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			var resp models.Error
			if rr.Code != http.StatusOK {
				err := json.NewDecoder(rr.Body).Decode(&resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Message, resp.Message)
			}
		})
	}
}

func generateMockToken(role string) string {
	return role
}

func TestAddProductToReceptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    func(*MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное добавление товара в приёмку",
			requestBody: map[string]interface{}{
				"type":  "электроника",
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("AddProductToActiveReception", mock.Anything, "электроника", uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Product{
						ID:          uuid.New(),
						Type:        "электроника",
						DateTime:    time.Now(),
						ReceptionID: uuid.New(),
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"type":"электроника"`,
		}, {
			name: "Ошибка при добавлении товара из-за отсутствия активной приёмки",
			requestBody: map[string]interface{}{
				"type":  "одежда",
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("AddProductToActiveReception", mock.Anything, "одежда", uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Product{}, apperrors.ErrNoActiveReception)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"нет активной приёмки для данного ПВЗ"`,
		}, {
			name: "Ошибка при добавлении товара из-за отсутствия активной приёмки",
			requestBody: map[string]interface{}{
				"type":  "одежда",
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("AddProductToActiveReception", mock.Anything, "одежда", uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Product{}, apperrors.ErrNoActiveReception)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"нет активной приёмки для данного ПВЗ"`,
		},
		{
			name:           "Некорректный JSON",
			requestBody:    `{"invalid_json"`,
			mockService:    func(m *MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Неверный формат запроса"`,
		},
		{
			name: "Ошибка при добавлении товара из-за внутренней ошибки сервера",
			requestBody: map[string]interface{}{
				"type":  "обувь",
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("AddProductToActiveReception", mock.Anything, "обувь", uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Product{}, errors.New("внутренняя ошибка сервера"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"Ошибка сервера: внутренняя ошибка сервера"`,
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

			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.addProductToReceptionHandler(rec, req)

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

func TestCreateReceptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockService    func(*MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное создание приёмки",
			requestBody: map[string]string{
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("CreateReception", mock.Anything, uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Reception{
						ID:       uuid.New(),
						DateTime: time.Now(),
						PVZID:    uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885"),
						Status:   "in_progress",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"status":"in_progress"`,
		},
		{
			name: "Ошибка сервиса при создании приёмки",
			requestBody: map[string]string{
				"pvzId": "86a4c84c-9719-419c-8449-f03267a2c885",
			},
			mockService: func(m *MockService) {
				m.On("CreateReception", mock.Anything, uuid.MustParse("86a4c84c-9719-419c-8449-f03267a2c885")).
					Return(models.Reception{}, errors.New("ошибка сервиса"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"Внутренняя ошибка сервера"`,
		},
		{
			name:           "Некорректный JSON",
			requestBody:    `{"invalid_json"`,
			mockService:    func(m *MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Невалидный JSON"`,
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

			req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.createReceptionHandler(rec, req)

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
