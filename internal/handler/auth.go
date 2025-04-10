package handler

import (
	"encoding/json"
	"github.com/gookit/slog"
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

	if req.Role != "client" && req.Role != "moderator" && req.Role != "pvz_employee" {
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
