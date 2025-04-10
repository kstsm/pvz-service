package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kstsm/pvz-service/config"
	"log"
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

	log.Printf("Генерация токена для пользователя с ролью: %s", role)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Ошибка при подписании токена: %v\n", err)
		return "", err
	}

	log.Printf("Токен успешно сгенерирован для пользователя с ролью: %s", role)
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Ошибка: неверный метод подписи: %v\n", token.Header["alg"])
			return nil, errors.New("неверный метод подписи")
		}
		return secretKey, nil
	})
	if err != nil {
		log.Printf("Ошибка при валидации токена: %v\n", err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Неверный токен или повреждённые данные")
		return "", errors.New("неверный токен")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		log.Println("Поле exp отсутствует или неверного типа")
		return "", errors.New("поле exp отсутствует или неверного типа")
	}
	if time.Now().Unix() > int64(expFloat) {
		log.Println("Токен истёк")
		return "", errors.New("токен истёк")
	}

	role, ok := claims["role"].(string)
	if !ok {
		log.Println("Поле role отсутствует или неверного типа")
		return "", errors.New("поле role отсутствует или неверного типа")
	}

	log.Printf("Токен успешно валидирован, роль пользователя: %s", role)
	return role, nil
}
