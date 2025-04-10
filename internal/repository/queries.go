package repository

const (
	QueryCreatePVZ = `
	INSERT INTO pvz (city)
	VALUES ($1)
	RETURNING id, registration_date, city;`

	QueryCreateReception = `
		WITH existing_reception AS (
            SELECT id, status
            FROM receptions
            WHERE pvz_id = $1 AND status != 'close'
            LIMIT 1
        )
        INSERT INTO receptions (pvz_id, status)
        SELECT $1, 'in_progress'
        WHERE NOT EXISTS (SELECT 1 FROM existing_reception)
        RETURNING id,date_time, pvz_id, status;`

	QueryAddProductToReception = `
		WITH existing_reception AS (
            SELECT id
            FROM receptions
            WHERE pvz_id = $1 AND status != 'close'
            ORDER BY date_time DESC
            LIMIT 1
        )
        INSERT INTO products (type, reception_id)
        SELECT $2, (SELECT id FROM existing_reception)
        WHERE EXISTS (SELECT 1 FROM existing_reception)
        RETURNING id, date_time,type, reception_id;`

	QueryGetLastOpenReception = `
		SELECT id, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1`

	QueryCloseReception = `
		UPDATE receptions
		SET status = 'close'
		WHERE id = $1
		RETURNING id, pvz_id, status, date_time`
)
