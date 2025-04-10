package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/kstsm/pvz-service/models"
	"net/http"
	"strings"
)

func (h Handler) closeLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "pvz_employee" {
		http.Error(w, "Доступ запрещен: только сотрудники ПВЗ могут закрывать приемки", http.StatusForbidden)
		return
	}

	pvzIDParam := chi.URLParam(r, "pvzId")
	pvzID, err := uuid.Parse(pvzIDParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный идентификатор ПВЗ: %v", err), http.StatusBadRequest)
		return
	}

	reception, err := h.service.CloseLastReception(r.Context(), pvzID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при закрытии приёмки: %v", err), http.StatusBadRequest)
		return
	}

	sendJSONResponse(w, http.StatusCreated, reception)
}

func (h Handler) deleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "pvz_employee" {
		writeErrorResponse(w, http.StatusForbidden, "Доступ запрещен: только сотрудники ПВЗ могут удалять товары")
		return
	}

	pvzIDStr := chi.URLParam(r, "pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "неверный формат pvzId")
		return
	}

	err = h.service.DeleteLastProductInReception(r.Context(), pvzID)
	if err != nil {
		if strings.Contains(err.Error(), "нет активной приёмки") {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("внутренняя ошибка: %s", err.Error()))
		return
	}

	sendJSONResponse(w, http.StatusOK, nil)
}

func (h Handler) createReceptionHandler(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "pvz_employee" {
		writeErrorResponse(w, http.StatusForbidden, "Доступ запрещен: только сотрудники ПВЗ могут создавать приемку")
		return
	}

	var req models.Reception
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	reception, err := h.service.CreateReception(r.Context(), req.PVZID)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSONResponse(w, http.StatusCreated, reception)
}

func (h Handler) addProductToReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	if req.Type == "" || req.PVZID == uuid.Nil {
		http.Error(w, `{"message": "тип товара и pvzId обязательны"}`, http.StatusBadRequest)
		return
	}

	product, err := h.service.AddProductToReception(r.Context(), req.Type, req.PVZID)
	if err != nil {
		if strings.Contains(err.Error(), "нет активной приёмки") {
			http.Error(w, fmt.Sprintf(`{"message": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf(`{"message": "ошибка сервера: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, product)
}
