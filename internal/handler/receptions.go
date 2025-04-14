package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
	"net/http"
)

func (h Handler) createReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Reception

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Некорректный JSON в теле запроса при создании приёмки", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Невалидный JSON")
		return
	}

	reception, err := h.service.CreateReception(r.Context(), req.PVZID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrReceptionAlreadyInProgress):
			slog.Info("Попытка создания приёмки при уже открытой приёмке", "pvzId", req.PVZID)
			writeErrorResponse(w, http.StatusBadRequest, "Невозможно создать приёмку: предыдущая не закрыта")
		default:
			slog.Error("Внутренняя ошибка при создании приёмки", "error", err, "pvzId", req.PVZID)
			writeErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	sendJSONResponse(w, http.StatusCreated, reception)
}

func (h Handler) addProductToReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Ошибка декодирования JSON при добавлении товара", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	if req.Type == "" || req.PVZID == uuid.Nil {
		slog.Warn("Отсутствуют обязательные поля в запросе на добавление товара", "req", req)
		writeErrorResponse(w, http.StatusBadRequest, "Тип товара и PVZ ID обязательны")
		return
	}

	if !isValidProduct(req.Type) {
		slog.Error("Недопустимый продукт", "role", req.Type)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимый продукт")
		return
	}

	product, err := h.service.AddProductToActiveReception(r.Context(), req.Type, req.PVZID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNoActiveReception):
			slog.Warn("Ошибка при добавлении товара: нет активной приёмки", "req", req, "error", err)
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			slog.Warn("Ошибка при добавлении товара в приёмку", "req", req, "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка сервера: "+err.Error())
		}
		return
	}

	sendJSONResponse(w, http.StatusCreated, product)
}

func (h Handler) deleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzIDStr := chi.URLParam(r, "pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		slog.Warn("Некорректный UUID ПВЗ при удалении товара", "pvzId", pvzIDStr, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат идентификатора ПВЗ")
		return
	}

	err = h.service.DeleteLastProductInReception(r.Context(), pvzID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNoActiveReception):
			slog.Warn("Нет активной приёмки для удаления товара", "pvzId", pvzID, "error", err)
			writeErrorResponse(w, http.StatusBadRequest, "Нет активной приёмки для удаления товара")
		case errors.Is(err, apperrors.ErrNoProductToDelete):
			slog.Warn("Нет товаров для удаления в активной приёмке", "pvzId", pvzID, "error", err)
			writeErrorResponse(w, http.StatusBadRequest, "Нет товаров для удаления в активной приёмке")
		default:
			slog.Error("Ошибка при удалении товара из приёмки", "pvzId", pvzID, "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Внутренняя ошибка: %s", err.Error()))
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, nil)
}

func (h Handler) closeLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzIDParam := chi.URLParam(r, "pvzId")
	pvzID, err := uuid.Parse(pvzIDParam)
	if err != nil {
		slog.Warn("Некорректный UUID ПВЗ при закрытии приёмки", "pvzId", pvzIDParam, "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный идентификатор ПВЗ")
		return
	}

	reception, err := h.service.CloseLastReception(r.Context(), pvzID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrReceptionAlreadyClosed):
			slog.Warn("Попытка закрыть уже закрытую или не найденную приемку", "pvzId", pvzID, "error", err)
			writeErrorResponse(w, http.StatusBadRequest, "Приемка уже закрыта или не найдена")
		default:
			slog.Error("Ошибка при закрытии приемки", "pvzId", pvzID, "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка сервера: "+err.Error())
		}

		return
	}

	sendJSONResponse(w, http.StatusOK, reception)
}
