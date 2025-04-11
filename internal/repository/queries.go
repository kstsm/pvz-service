package repository

const (
	queryCreatePVZ = `
		INSERT INTO pvz (city)
		VALUES ($1)
		RETURNING id, registration_date, city;`

	queryCreateReception = `
    	INSERT INTO receptions (pvz_id, status)
		SELECT $1, 'in_progress'
		WHERE NOT EXISTS (
    	SELECT 1 FROM receptions WHERE pvz_id = $1 AND status != 'close'
		)
		RETURNING id, date_time, pvz_id, status;`

	queryGetLastOpenReception = `
		SELECT id, pvz_id, status
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
		LIMIT 1`

	queryCloseReception = `
		UPDATE receptions
		SET status = 'close'
		WHERE id = $1
		RETURNING id, pvz_id, status, date_time`

	queryCheckEmailExists = `
		SELECT EXISTS(
		SELECT 1 
		FROM users 
		WHERE email = $1)`

	queryCreateUser = `
		INSERT INTO users(email, password, role) 
		VALUES($1, $2, $3) RETURNING id`

	queryGetRoleByEmail = `
		SELECT password, role 
		FROM users 
		WHERE email = $1`

	queryGetActiveReception = `
		SELECT id
		FROM receptions
		WHERE pvz_id = $1 AND status != 'close'
		ORDER BY date_time DESC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`
	queryInsertProduct = `
		INSERT INTO products (id, type, reception_id)
		VALUES ($1, $2, $3)
		RETURNING id, date_time, type, reception_id
	`

	getLastProductQuery = `
	SELECT p.id
	FROM products p
	WHERE p.reception_id = (
		SELECT r.id
		FROM receptions r
		WHERE r.pvz_id = $1 AND r.status != 'close'
		LIMIT 1
	) 
	ORDER BY p.date_time DESC
	LIMIT 1
	FOR UPDATE SKIP LOCKED;
`
)
