package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		tokenHeader    string
		validateToken  func(string) (string, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Отсутствует токен",
			tokenHeader:    "",
			validateToken:  nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"Отсутствует токен авторизации"}`,
		},
		{
			name:        "Неверный токен",
			tokenHeader: "Bearer invalid-token",
			validateToken: func(token string) (string, error) {
				return "", assert.AnError
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"errors":"Неверный или просроченный токен"}`,
		},
		{
			name:        "Валидный токен",
			tokenHeader: "Bearer valid-token",
			validateToken: func(token string) (string, error) {
				return "admin", nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.tokenHeader != "" {
				req.Header.Set("Authorization", tt.tokenHeader)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})

			middleware := AuthMiddleware(tt.validateToken)(handler)
			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if rr.Body.Len() > 0 {
				var body map[string]string
				if err := json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&body); err == nil {
					assert.Equal(t, tt.expectedBody, `{"errors":"`+body["errors"]+`"}`)
				} else {
					assert.Equal(t, "ok", rr.Body.String())
				}
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		role           any
		allowedRoles   []string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Роль отсутствует",
			role:           nil,
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"errors":"Роль не найдена в контексте"}`,
		},
		{
			name:           "Роль не разрешена",
			role:           "user",
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"errors":"Недостаточно прав доступа"}`,
		},
		{
			name:           "Роль разрешена",
			role:           "admin",
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			ctx := context.WithValue(req.Context(), "role", tt.role)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			})

			middleware := RequireRole(tt.allowedRoles...)(handler)
			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if rr.Body.Len() > 0 {
				var body map[string]string
				if err := json.NewDecoder(bytes.NewReader(rr.Body.Bytes())).Decode(&body); err == nil {
					assert.Equal(t, tt.expectedBody, `{"errors":"`+body["errors"]+`"}`)
				} else {
					assert.Equal(t, "ok", rr.Body.String())
				}
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		expected   string
	}{
		{
			name:       "Корректный токен",
			authHeader: "Bearer my-secret-token",
			expected:   "my-secret-token",
		},
		{
			name:       "Пустой заголовок",
			authHeader: "",
			expected:   "",
		},
		{
			name:       "Неверный формат - только токен без Bearer",
			authHeader: "my-secret-token",
			expected:   "",
		},
		{
			name:       "Неверный формат - Bearer и ничего",
			authHeader: "Bearer",
			expected:   "",
		},
		{
			name:       "Неверный формат - три части",
			authHeader: "Bearer token extra",
			expected:   "",
		},
		{
			name:       "Неверный префикс",
			authHeader: "Token my-token",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			token := ExtractToken(req)
			assert.Equal(t, tt.expected, token)
		})
	}
}

func TestSendJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	sendJSONError(rr, http.StatusUnauthorized, "Ошибка авторизации")

	resp := rr.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body struct {
		Errors string `json:"errors"`
	}
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "Ошибка авторизации", body.Errors)
}
