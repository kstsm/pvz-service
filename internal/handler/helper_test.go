package handler

import (
	"encoding/json"
	"github.com/kstsm/pvz-service/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsValidProduct(t *testing.T) {
	tests := []struct {
		product  string
		expected bool
	}{
		{"электроника", true},
		{"одежда", true},
		{"обувь", true},
		{"мебель", false},
		{"косметика", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.product, func(t *testing.T) {
			actual := isValidProduct(test.product)
			if actual != test.expected {
				t.Errorf("isValidProduct(%s) = %v; expected %v", test.product, actual, test.expected)
			}
		})
	}
}

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"client", true},
		{"moderator", true},
		{"employee", true},
		{"admin", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidRole(tt.role))
		})
	}
}

func TestIsValidCity(t *testing.T) {
	tests := []struct {
		city     string
		expected bool
	}{
		{"Москва", true},
		{"Санкт-Петербург", true},
		{"Казань", true},
		{"Новосибирск", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.city, func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidCity(tt.city))
		})
	}
}

func TestSendJSONResponse(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]string{"message": "ok"}
	sendJSONResponse(w, http.StatusOK, data)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body map[string]string
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", body["message"])
}

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	writeErrorResponse(w, http.StatusBadRequest, "ошибка валидации")

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body models.Error
	err := json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "ошибка валидации", body.Message)
}

func TestParsePVZFilterParams(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedParams models.PVZFilterParams
		expectedError  bool
	}{
		{
			name: "Valid params",
			queryParams: map[string]string{
				"startDate": "2025-04-01T00:00:00Z",
				"endDate":   "2025-04-10T00:00:00Z",
				"page":      "2",
				"limit":     "20",
			},
			expectedParams: models.PVZFilterParams{
				StartDate: parseTime("2025-04-01T00:00:00Z"),
				EndDate:   parseTime("2025-04-10T00:00:00Z"),
				Page:      2,
				Limit:     20,
			},
			expectedError: false,
		},
		{
			name: "Invalid startDate format",
			queryParams: map[string]string{
				"startDate": "invalid-date",
			},
			expectedParams: models.PVZFilterParams{
				Page:  1,
				Limit: 10,
			},
			expectedError: true,
		},
		{
			name: "Invalid page value",
			queryParams: map[string]string{
				"page": "0",
			},
			expectedParams: models.PVZFilterParams{
				Page:  1,
				Limit: 10,
			},
			expectedError: false,
		},
		{
			name: "Limit out of range",
			queryParams: map[string]string{
				"limit": "50",
			},
			expectedParams: models.PVZFilterParams{
				Page:  1,
				Limit: 10,
			},
			expectedError: false,
		},
		{
			name: "Missing startDate and endDate",
			queryParams: map[string]string{
				"page":  "3",
				"limit": "5",
			},
			expectedParams: models.PVZFilterParams{
				Page:  3,
				Limit: 5,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			actualParams, err := parsePVZFilterParams(req)

			if tt.expectedError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("did not expect error, got: %v", err)
			}

			if !equalParams(actualParams, tt.expectedParams) {
				t.Errorf("expected %v, got %v", tt.expectedParams, actualParams)
			}
		})
	}
}

func parseTime(timeStr string) *time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	return &parsedTime
}

func equalParams(a, b models.PVZFilterParams) bool {
	if a.Page != b.Page || a.Limit != b.Limit {
		return false
	}
	if a.StartDate != nil && b.StartDate != nil && !a.StartDate.Equal(*b.StartDate) {
		return false
	}
	if a.EndDate != nil && b.EndDate != nil && !a.EndDate.Equal(*b.EndDate) {
		return false
	}
	return true
}
