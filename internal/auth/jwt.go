package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/config"
	"time"
)

const tokenExpiry = time.Hour * 24

func GenerateToken(role string) (string, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		slog.Error("Ошибка при подписании токена: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("Ошибка: неверный метод подписи", "alg", token.Header["alg"])
			return nil, errors.New("неверный метод подписи")
		}
		return secretKey, nil
	})
	if err != nil {
		slog.Error("Ошибка при валидации токена", "error", err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Error("Неверный токен или повреждённые данные")
		return "", errors.New("неверный токен")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		slog.Error("Поле exp отсутствует или неверного типа")
		return "", errors.New("поле exp отсутствует или неверного типа")
	}
	if time.Now().Unix() > int64(expFloat) {
		slog.Error("Токен истёк")
		return "", errors.New("токен истёк")
	}

	role, ok := claims["role"].(string)
	if !ok {
		slog.Error("Поле role отсутствует или неверного типа")
		return "", errors.New("поле role отсутствует или неверного типа")
	}

	return role, nil
}
