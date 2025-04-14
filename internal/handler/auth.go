package handler

import (
	"encoding/json"
	"errors"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/internal/auth"
	"github.com/kstsm/pvz-service/models"
	"net/http"
)

func (h Handler) dummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DummyLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Недопустимое тело запроса", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимое тело запроса")
		return
	}

	if !isValidRole(req.Role) {
		slog.Error("Недопустимая роль", "role", req.Role)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимая роль")
		return
	}

	token, err := auth.GenerateToken(req.Role)
	if err != nil {
		slog.Error("Ошибка генерации токена", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Пользователь не найден")
		return
	}

	sendJSONResponse(w, http.StatusOK, token)
}

func (h Handler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Ошибка при декодировании тела запроса", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимое тело запроса")
		return
	}

	if !isValidRole(req.Role) {
		slog.Error("Недопустимая роль", "role", req.Role)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимая роль")
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrEmailAlreadyExists):
			writeErrorResponse(w, http.StatusConflict, "Пользователь с таким email уже существует")
		default:
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при регистрации пользователя")
		}
		slog.Error("Ошибка регистрации пользователя", "email", req.Email, "error", err)
		return
	}

	sendJSONResponse(w, http.StatusCreated, user)
}

func (h Handler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Ошибка при декодировании тела запроса", "error", err)
		writeErrorResponse(w, http.StatusBadRequest, "Недопустимое тело запроса")
		return
	}

	userRole, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			writeErrorResponse(w, http.StatusUnauthorized, "Неверный email или пароль")
		case errors.Is(err, apperrors.ErrEmailNotFound):
			writeErrorResponse(w, http.StatusNotFound, "Пользователь с таким email не найден")
		default:
			slog.Error("Ошибка авторизации пользователя", "email", req.Email, "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при авторизации пользователя")
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, userRole)
}
