package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        models.UserRegisterReq
		mockService    func() *MockService
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Успешная_регистрация_-_client",
			reqBody: models.UserRegisterReq{
				Email:    "client@example.com",
				Password: "password123",
				Role:     "client",
			},
			mockService: func() *MockService {
				return &MockService{}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Успешная_регистрация_-_moderator",
			reqBody: models.UserRegisterReq{
				Email:    "moderator@example.com",
				Password: "password123",
				Role:     "moderator",
			},
			mockService: func() *MockService {
				return &MockService{}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Успешная_регистрация_-_employee",
			reqBody: models.UserRegisterReq{
				Email:    "employee@example.com",
				Password: "password123",
				Role:     "employee",
			},
			mockService: func() *MockService {
				return &MockService{}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "Недопустимая_роль",
			reqBody: models.UserRegisterReq{
				Email:    "admin@example.com",
				Password: "password123",
				Role:     "admin",
			},
			mockService: func() *MockService {
				return &MockService{}
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Недопустимая роль",
		},

		{
			name: "Пользователь_с_таким_email_уже_существует",
			reqBody: models.UserRegisterReq{
				Email:    "existing@example.com",
				Password: "password123",
				Role:     "client",
			},
			mockService: func() *MockService {
				return &MockService{}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Пользователь с таким email уже существует",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBodyBytes, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatal("Не удалось маршалить тело запроса", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBodyBytes))
			rec := httptest.NewRecorder()

			handler := Handler{
				service: tt.mockService(),
			}

			handler.registerUserHandler(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedError != "" {
				var errResp models.Error
				if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
					t.Fatal("Не удалось декодировать ошибку ответа", err)
				}
				assert.Equal(t, tt.expectedError, errResp.Message)
			} else {
				var userResp models.UserRegisterResp
				if err := json.NewDecoder(res.Body).Decode(&userResp); err != nil {
					t.Fatal("Не удалось декодировать ответ", err)
				}
				assert.NotEmpty(t, userResp.ID)
				assert.Equal(t, tt.reqBody.Email, userResp.Email)
				assert.Equal(t, tt.reqBody.Role, userResp.Role)
			}
		})
	}
}

func TestDummyLoginHandler(t *testing.T) {
	type args struct {
		Role string `json:"role"`
	}

	tests := []struct {
		name       string
		body       args
		wantStatus int
		wantRole   string
	}{
		{
			name:       "Успешный вход - client",
			body:       args{Role: "client"},
			wantStatus: http.StatusOK,
			wantRole:   "client",
		},
		{
			name:       "Успешный вход - moderator",
			body:       args{Role: "moderator"},
			wantStatus: http.StatusOK,
			wantRole:   "moderator",
		},
		{
			name:       "Успешный вход - employee",
			body:       args{Role: "employee"},
			wantStatus: http.StatusOK,
			wantRole:   "employee",
		},
		{
			name:       "Недопустимая роль",
			body:       args{Role: "admin"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Пустое тело запроса",
			body:       args{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error

			mockService := new(MockService)
			handler := Handler{service: mockService}

			if tt.name != "Пустое тело запроса" {
				body, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/dummy/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.dummyLoginHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			require.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				var buf bytes.Buffer
				_, err := buf.ReadFrom(resp.Body)
				require.NoError(t, err)

				tokenString := strings.Trim(buf.String(), "\"\n")

				token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
				require.NoError(t, err)

				claims, ok := token.Claims.(jwt.MapClaims)
				require.True(t, ok)
				require.Equal(t, tt.wantRole, claims["role"])
			}
		})
	}
}

func TestLoginUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		loginReq       models.UserLoginReq
		mockResponse   string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Успешная авторизация",
			loginReq:       models.UserLoginReq{Email: "user@example.com", Password: "password123"},
			mockResponse:   "user_role_token",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"user_role_token"`,
		},
		{
			name:           "Неверные учетные данные",
			loginReq:       models.UserLoginReq{Email: "user@example.com", Password: "wrongpassword"},
			mockResponse:   "",
			mockError:      apperrors.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Неверный email или пароль"}`,
		},
		{
			name:           "Пользователь не найден",
			loginReq:       models.UserLoginReq{Email: "nonexistent@example.com", Password: "password123"},
			mockResponse:   "",
			mockError:      apperrors.ErrEmailNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Пользователь с таким email не найден"}`,
		},
		{
			name:           "Ошибка сервера",
			loginReq:       models.UserLoginReq{Email: "user@example.com", Password: "password123"},
			mockResponse:   "",
			mockError:      errors.New("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Ошибка при авторизации пользователя"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockService)
			mockService.On("LoginUser", mock.Anything, tt.loginReq).Return(tt.mockResponse, tt.mockError)

			handler := Handler{service: mockService}

			reqBody, err := json.Marshal(tt.loginReq)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			handler.loginUserHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
