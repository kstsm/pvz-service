package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kstsm/pvz-service/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func isValidProduct(product string) bool {
	allowedProduct := map[string]bool{
		"электроника": true,
		"одежда":      true,
		"обувь":       true,
	}
	return allowedProduct[product]
}

func isValidRole(role string) bool {
	allowedRoles := map[string]bool{
		"client":    true,
		"moderator": true,
		"employee":  true,
	}
	return allowedRoles[role]
}

func isValidCity(city string) bool {
	allowedCity := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}
	return allowedCity[city]
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.Error{Message: message})
}

func parsePVZFilterParams(r *http.Request) (models.PVZFilterParams, error) {
	params := models.PVZFilterParams{
		Page:  1,
		Limit: 10,
	}

	if startDate := r.URL.Query().Get("startDate"); startDate != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return params, fmt.Errorf("неверный формат даты начала: %v", err)
		}
		params.StartDate = &parsedStartDate
	}

	if endDate := r.URL.Query().Get("endDate"); endDate != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return params, fmt.Errorf("неверный формат даты конца: %v", err)
		}
		params.EndDate = &parsedEndDate
	}

	if page := r.URL.Query().Get("page"); page != "" {
		parsedPage, err := strconv.Atoi(page)
		if err == nil && parsedPage > 0 {
			params.Page = parsedPage
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil && parsedLimit > 0 && parsedLimit <= 30 {
			params.Limit = parsedLimit
		}
	}

	return params, nil
}
