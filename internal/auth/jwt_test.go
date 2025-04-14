package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kstsm/pvz-service/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testSecret = "test-secret"

func TestGenerateToken(t *testing.T) {
	config.Config.JWT.JWTSecret = testSecret

	t.Run("успешная генерация токена", func(t *testing.T) {
		token, err := GenerateToken("employee")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(testSecret), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsed.Valid)

		claims, ok := parsed.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "employee", claims["role"])
	})
}

func TestValidateToken(t *testing.T) {
	config.Config.JWT.JWTSecret = "my_secret"

	t.Run("валидный токен", func(t *testing.T) {
		token, err := GenerateToken("admin")
		assert.NoError(t, err)

		role, err := ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "admin", role)
	})

	t.Run("невалидный токен", func(t *testing.T) {
		_, err := ValidateToken("this.is.not.valid")
		assert.Error(t, err)
	})

	t.Run("токен без role", func(t *testing.T) {
		claims := jwt.MapClaims{
			"exp": time.Now().Add(10 * time.Minute).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(config.Config.JWT.JWTSecret))

		_, err := ValidateToken(tokenStr)
		assert.EqualError(t, err, "поле role отсутствует или неверного типа")
	})

	t.Run("неподдерживаемый метод подписи", func(t *testing.T) {
		claims := jwt.MapClaims{
			"role": "admin",
			"exp":  time.Now().Add(10 * time.Minute).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		tokenStr, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

		_, err := ValidateToken(tokenStr)
		assert.EqualError(t, err, "неверный метод подписи")
	})
}
