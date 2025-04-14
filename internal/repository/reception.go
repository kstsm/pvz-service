package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/pvz-service/internal/apperrors"
	"github.com/kstsm/pvz-service/models"
)

func (r Repository) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	var reception models.Reception

	err := r.conn.QueryRow(ctx, queryCreateReception, pvzID).
		Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Info("Приёмка не создана: существует незакрытая приемка", "pvzId", pvzID)
			return models.Reception{}, apperrors.ErrReceptionAlreadyInProgress
		}

		slog.Error("Ошибка при заведении приемки", "error", err, "pvzId", pvzID)
		return models.Reception{}, fmt.Errorf("не удалось создать приёмку: %w", err)
	}

	return reception, nil
}

func (r Repository) AddProductToActiveReception(ctx context.Context, productType string, pvzID uuid.UUID) (models.Product, error) {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return models.Product{}, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback(ctx)

	var receptionID uuid.UUID
	err = tx.QueryRow(ctx, queryGetActiveReception, pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, apperrors.ErrNoActiveReception
		}
		return models.Product{}, fmt.Errorf("ошибка при получении активной приёмки: %w", err)
	}

	newID := uuid.New()
	var product models.Product
	err = tx.QueryRow(ctx, queryInsertProduct, newID, productType, receptionID).
		Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if err != nil {
		return models.Product{}, fmt.Errorf("ошибка при добавлении товара: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return models.Product{}, fmt.Errorf("не удалось завершить транзакцию: %w", err)
	}

	return product, nil
}

func (r Repository) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Ошибка при начале транзакции", "pvzId", pvzID, "error", err)
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback(ctx)

	var productID uuid.UUID

	var receptionExists bool
	err = tx.QueryRow(ctx, checkActiveReceptionQuery, pvzID).Scan(&receptionExists)
	if err != nil {
		slog.Error("Ошибка при проверке активной приемки", "pvzId", pvzID, "error", err)
		return fmt.Errorf("ошибка при проверке активной приемки: %w", err)
	}

	if !receptionExists {
		slog.Warn("Нет активной приемки", "pvzId", pvzID)
		return apperrors.ErrNoActiveReception
	}

	err = tx.QueryRow(ctx, getLastProductQuery, pvzID).Scan(&productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Warn("Нет товаров для удаления", "pvzId", pvzID)
			return apperrors.ErrNoProductToDelete
		}

		slog.Error("Ошибка при получении последнего товара", "pvzId", pvzID, "error", err)
		return fmt.Errorf("ошибка при получении последнего товара: %w", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		slog.Error("Ошибка при удалении товара", "pvzId", pvzID, "productID", productID, "error", err)
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("Ошибка при фиксации транзакции", "pvzId", pvzID, "error", err)
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}

func (r Repository) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	var reception models.Reception

	err := r.conn.QueryRow(ctx, queryGetLastOpenReception, pvzID).Scan(&reception.ID, &reception.PVZID, &reception.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Reception{}, apperrors.ErrReceptionAlreadyClosed
		}
		return models.Reception{}, fmt.Errorf("не удалось получить последнюю открытую приёмку для ПВЗ с ID %v: %w", pvzID, err)
	}

	err = r.conn.QueryRow(ctx, queryCloseReception, reception.ID).Scan(&reception.ID, &reception.PVZID, &reception.Status, &reception.DateTime)
	if err != nil {
		return models.Reception{}, fmt.Errorf("не удалось закрыть приемку с ID %v для ПВЗ с ID %v: %w", reception.ID, pvzID, err)
	}

	return reception, nil
}
