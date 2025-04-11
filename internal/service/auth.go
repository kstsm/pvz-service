package service

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/internal/auth"
	"github.com/kstsm/pvz-service/models"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) RegisterUser(ctx context.Context, req models.UserRegisterReq) (models.UserRegisterResp, error) {
	exists, err := s.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		slog.Error("Ошибка при проверке email", "email", req.Email, "error", err)
		return models.UserRegisterResp{}, fmt.Errorf("не удалось проверить email: %w", err)
	}
	if exists {
		slog.Error("Пользователь с таким email уже существует", "email", req.Email)
		return models.UserRegisterResp{}, apperrors.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Ошибка хеширования пароля", "email", req.Email, "error", err)
		return models.UserRegisterResp{}, fmt.Errorf("не удалось хешировать пароль: %w", err)
	}
	req.Password = string(hashedPassword)

	userID, err := s.repo.CreateUser(ctx, req)
	if err != nil {
		slog.Error("Ошибка регистрации пользователя", "email", req.Email, "error", err)
		return models.UserRegisterResp{}, fmt.Errorf("не удалось зарегистрировать пользователя: %w", err)
	}

	userResp := models.UserRegisterResp{
		ID:    userID,
		Email: req.Email,
		Role:  req.Role,
	}

	return userResp, nil
}

func (s Service) LoginUser(ctx context.Context, req models.UserLoginReq) (string, error) {
	hashedPassword, userRole, err := s.repo.GetRoleByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		slog.Warn("Неверный пароль", "email", req.Email)
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(userRole)
	if err != nil {
		slog.Error("Ошибка генерации токена", "error", err)
		return "", fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return token, nil
}
