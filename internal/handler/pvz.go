package handler

import (
	"encoding/json"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/models"
	"net/http"
)

func (h Handler) getListPVZ(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok || (role != "moderator" && role != "pvz_employee") {
		writeErrorResponse(w, http.StatusForbidden, "Доступ запрещен: только модераторы и ПВЗ работники могут выполнить данное действие")
		slog.Warn("Запрет доступа к списку ПВЗ", "role", role)
		return
	}

	params, err := parsePVZFilterParams(r)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		slog.Warn("Невалидные параметры запроса", "error", err)
		return
	}

	pvzList, err := h.service.GetPVZList(r.Context(), params)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось получить список ПВЗ")
		slog.Error("Ошибка получения списка ПВЗ", "error", err)
		return
	}

	sendJSONResponse(w, http.StatusOK, pvzList)
}

func (h Handler) createPVZHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	if role != "moderator" {
		writeErrorResponse(w, http.StatusForbidden, "Доступ запрещен: только модераторы могут создавать ПВЗ")
		return
	}

	var req models.PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	if req.City != "Москва" && req.City != "Санкт-Петербург" && req.City != "Казань" {
		writeErrorResponse(w, http.StatusBadRequest, "Можно создать ПВЗ только в городах: Москва, Санкт-Петербург, Казань")
		return
	}

	pvz, err := h.service.CreatePVZ(r.Context(), req.City)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось создать ПВЗ")
		return
	}

	sendJSONResponse(w, http.StatusOK, pvz)
}
