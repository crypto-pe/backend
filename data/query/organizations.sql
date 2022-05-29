-- id	uuid [uuid_generate_v4()]	
-- name	character varying(60)	
-- created_at	timestamp	
-- owner_address	bytea	
-- token	bytea

-- name: CreateOrganization :one
INSERT INTO organizations(name, owner_address, token)
VALUES ($1, $2, $3)
RETURNING *;

-- wtf?
-- wdym?
-- it expects uuid?
-- nv
-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1;

-- name: UpdateOrganization :one
UPDATE organizations
SET
  name = $2,
  token = $3
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: GetAllOrganizations :many
SELECT * FROM organizations WHERE id IN (SELECT organization_id FROM organization_members WHERE member_address=$1);
