package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
)

func (r Repository) CreateUser(ctx context.Context, user models.UserRegisterReq) (uuid.UUID, error) {
	var id uuid.UUID
	var pgError *pgconn.PgError

	err := r.conn.QueryRow(ctx, queryCreateUser, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && pgError.ConstraintName == "users_email_key" {
				return uuid.Nil, apperrors.ErrEmailAlreadyExists
			}
			return uuid.Nil, fmt.Errorf("не удалось зарегистрировать пользователя: %w", err)
		}
	}
	return id, nil
}

func (r Repository) GetRoleByEmail(ctx context.Context, email string) (string, string, error) {
	var hashedPassword, userRole string
	err := r.conn.QueryRow(ctx, queryGetRoleByEmail, email).Scan(&hashedPassword, &userRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", apperrors.ErrEmailNotFound
		}
		return "", "", fmt.Errorf("ошибка при получении данных пользователя по email: %w", err)
	}
	return hashedPassword, userRole, nil
}
