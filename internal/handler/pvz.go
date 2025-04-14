package handler

import (
	"encoding/json"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/models"
	"net/http"
)

func (h Handler) createPVZHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Ошибка декодирования JSON при создании ПВЗ", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	if !isValidCity(req.City) {
		slog.Warn("Попытка создать ПВЗ в недопустимом городе", "city", req.City)
		writeErrorResponse(w, http.StatusBadRequest, "Данный город пока недоступен")
		return
	}

	pvz, err := h.service.CreatePVZ(r.Context(), req.City)
	if err != nil {
		slog.Error("Ошибка при создании ПВЗ", "city", req.City, "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось создать ПВЗ")
		return
	}

	sendJSONResponse(w, http.StatusOK, pvz)
}

func (h Handler) getListPVZ(w http.ResponseWriter, r *http.Request) {
	params, err := parsePVZFilterParams(r)
	if err != nil {
		slog.Warn("Невалидные параметры запроса для фильтрации ПВЗ", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Невалидные параметры запроса")
		return
	}
	pvzList, err := h.service.GetPVZList(r.Context(), params)
	if err != nil {
		slog.Error("Ошибка при получении списка ПВЗ", "error", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось получить список ПВЗ")
		return
	}

	sendJSONResponse(w, http.StatusOK, pvzList)
}
