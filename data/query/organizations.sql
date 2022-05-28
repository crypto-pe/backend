-- id	uuid [uuid_generate_v4()]	
-- name	character varying(60)	
-- created_at	timestamp	
-- owner_address	bytea	
-- token	bytea

-- name: CreateOrganization :one
INSERT INTO organizations(name, owner_address, token)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1;

-- name: UpdateOrganization :one
UPDATE organizations
SET
  name = $2,
  owner_address = $3,
  token = $4
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;
