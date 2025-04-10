package middleware

import (
	"context"
	"encoding/json"
	"github.com/gookit/slog"
	"net/http"
	"strings"
)

func RoleMiddleware(validateToken func(string) (string, error), allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ExtractToken(r)
			if token == "" {
				sendJSONError(w, http.StatusUnauthorized, "Отсутствует токен авторизации")
				return
			}

			role, err := validateToken(token)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "Неверный или просроченный токен")
				return
			}

			allowed := false
			for _, ar := range allowedRoles {
				if role == ar {
					allowed = true
					break
				}
			}
			if !allowed {
				sendJSONError(w, http.StatusForbidden, "Недостаточно прав доступа")
				return
			}

			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		slog.Warn("Заголовок Authorization отсутствует", "path", r.URL.Path)
		return ""
	}
	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		slog.Warn("Некорректный формат токена", "path", r.URL.Path, "token", bearerToken)
		return ""
	}

	return parts[1]
}

func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := struct {
		Errors string `json:"errors"`
	}{
		Errors: message,
	}

	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		return
	}
}
