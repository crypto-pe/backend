-- address	bytea	
-- name	character varying(60)	
-- created_at	timestamp NULL [now()]	
-- email	domain_email	
-- admin	boolean NULL [false]

-- name: CreateUser :one
INSERT INTO accounts(address, name, email, admin)
VALUES ($1, $2, $3, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM accounts WHERE address = $1;

-- name: UpdateUser :one
UPDATE accounts 
SET
  name = $2,
  email = $3,
  admin = $4
WHERE address = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM accounts
WHERE address = $1;



