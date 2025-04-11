package middleware

import (
	"context"
	"encoding/json"
	"github.com/gookit/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(validateToken func(string) (string, error)) func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value("role")
			role, ok := roleVal.(string)
			if !ok || role == "" {
				sendJSONError(w, http.StatusForbidden, "Роль не найдена в контексте")
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			sendJSONError(w, http.StatusForbidden, "Недостаточно прав доступа")
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
