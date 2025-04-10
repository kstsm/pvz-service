package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/models"
)

func (r Repository) CreatePVZ(ctx context.Context, city string) (models.PVZ, error) {
	var pvz models.PVZ
	err := r.conn.QueryRow(ctx, QueryCreatePVZ, city).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		slog.Error("Ошибка при заведении ПВЗ", "error", err)
		return models.PVZ{}, fmt.Errorf("r.conn.QueryRow: %w", err)
	}

	return pvz, nil
}

func (r Repository) GetPVZList(ctx context.Context, params models.PVZFilterParams) ([]models.PVZWithReceptions, error) {
	var query string
	var args []interface{}
	query = `SELECT pvz.id, pvz.registration_date, pvz.city, receptions.id, receptions.date_time, receptions.status 
	          FROM pvz 
	          LEFT JOIN receptions ON pvz.id = receptions.pvz_id 
	          WHERE 1=1`

	if params.StartDate != nil {
		query += " AND receptions.date_time >= $1"
		args = append(args, *params.StartDate)
	}
	if params.EndDate != nil {
		query += " AND receptions.date_time <= $2"
		args = append(args, *params.EndDate)
	}

	query += " LIMIT $3 OFFSET $4"
	args = append(args, params.Limit, (params.Page-1)*params.Limit)

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pvzMap := make(map[uuid.UUID]*models.PVZWithReceptions)

	for rows.Next() {
		var pvz models.PVZ
		var reception models.Reception
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City, &reception.ID, &reception.DateTime, &reception.Status)
		if err != nil {
			return nil, err
		}

		if _, exists := pvzMap[pvz.ID]; !exists {
			pvzMap[pvz.ID] = &models.PVZWithReceptions{
				PVZ:        pvz,
				Receptions: []models.Reception{},
			}
		}

		pvzMap[pvz.ID].Receptions = append(pvzMap[pvz.ID].Receptions, reception)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var pvzList []models.PVZWithReceptions
	for _, pvzWithReceptions := range pvzMap {
		pvzList = append(pvzList, *pvzWithReceptions)
	}

	return pvzList, nil
}
