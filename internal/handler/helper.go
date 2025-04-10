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
	query := r.URL.Query()

	var resp models.PVZFilterParams

	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return resp, fmt.Errorf("invalid start_date format")
		}
		resp.StartDate = &startDate
	}

	endDateStr := r.URL.Query().Get("end_date")
	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return resp, fmt.Errorf("invalid end_date format")
		}
		resp.EndDate = &endDate
	}

	var err error

	pageStr := query.Get("page")
	resp.Page, err = strconv.Atoi(pageStr)
	if err != nil || resp.Page <= 0 {
		resp.Page = 1
	}

	limitStr := query.Get("limit")
	resp.Limit, err = strconv.Atoi(limitStr)
	if err != nil || resp.Limit <= 0 || resp.Limit > 30 {
		resp.Limit = 10
	}

	return resp, nil
}
