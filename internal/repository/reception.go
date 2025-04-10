package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5"
	"github.com/kstsm/pvz-service/models"
	"strings"
)

func (r Repository) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	var reception models.Reception

	err := r.conn.QueryRow(ctx, QueryGetLastOpenReception, pvzID).Scan(&reception.ID, &reception.PVZID, &reception.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Reception{}, fmt.Errorf("не найдена открытая приёмка для ПВЗ с ID %v", pvzID)
		}
		return models.Reception{}, fmt.Errorf("ошибка при получении приёмки: %w", err)
	}

	err = r.conn.QueryRow(ctx, QueryCloseReception, reception.ID).Scan(&reception.ID, &reception.PVZID, &reception.Status, &reception.DateTime)
	if err != nil {
		return models.Reception{}, fmt.Errorf("ошибка при закрытии приёмки: %w", err)
	}

	return reception, nil
}

func (r Repository) DeleteLastProductInReception(ctx context.Context, pvzID uuid.UUID) error {
	var productID uuid.UUID

	err := r.conn.QueryRow(ctx, `
		SELECT p.id
		FROM receptions r
		JOIN products p ON p.reception_id = r.id
		WHERE r.pvz_id = $1 AND r.status != 'close'
		ORDER BY p.date_time DESC
		LIMIT 1;
	`, pvzID).Scan(&productID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("нет активной приёмки или товары отсутствуют")
		}
		return fmt.Errorf("ошибка при получении последнего товара: %w", err)
	}

	_, err = r.conn.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товара: %w", err)
	}

	return nil
}

func (r Repository) CreateReception(ctx context.Context, pvzID uuid.UUID) (models.Reception, error) {
	var reception models.Reception
	err := r.conn.QueryRow(ctx, QueryCreateReception, pvzID).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		if strings.Contains(err.Error(), "NO_DATA_FOUND") {
			return models.Reception{}, fmt.Errorf("невозможно создать приемку: предыдущая приемка не была закрыта")
		}
		slog.Error("Ошибка при заведении приемки", "error", err)
		return models.Reception{}, fmt.Errorf("r.conn.QueryRow: %w", err)
	}
	return reception, nil
}

func (r Repository) GetActiveReceptionByPVZ(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error) {
	var receptionID uuid.UUID

	query := `
		SELECT id
		FROM receptions
		WHERE pvz_id = $1 AND status != 'close'
		ORDER BY date_time DESC
		LIMIT 1
	`

	err := r.conn.QueryRow(ctx, query, pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("нет активной приёмки для данного ПВЗ")
		}
		return uuid.Nil, fmt.Errorf("ошибка при получении приёмки: %w", err)
	}

	return receptionID, nil
}

func (r Repository) AddProduct(ctx context.Context, productType string, receptionID uuid.UUID) (models.Product, error) {
	var product models.Product

	query := `
		INSERT INTO products (id, type, reception_id)
		VALUES ($1, $2, $3)
		RETURNING id, date_time, type, reception_id
	`

	newID := uuid.New()

	err := r.conn.QueryRow(ctx, query, newID, productType, receptionID).
		Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID)
	if err != nil {
		return models.Product{}, fmt.Errorf("ошибка при добавлении товара: %w", err)
	}

	return product, nil
}
