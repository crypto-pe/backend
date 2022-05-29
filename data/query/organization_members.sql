-- organization_id	uuid	
-- member_address	bytea	
-- date_joined	timestamp	
-- role	character varying(200)	
-- is_admin	boolean NULL [false]	
-- salary	integer NULL [0]	

-- name: CreateOrganizationMember :one
INSERT INTO organization_members(
  organization_id, member_address, role, is_admin, salary
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetOrganizationMember :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND member_address = $2;

-- name: UpdateOrganizationMember :one
UPDATE organization_members
SET
  role = $3, 
  is_admin = $4,
  salary = $5
WHERE organization_id = $1 AND member_address = $2
RETURNING *;

-- name: DeleteOrganizationMember :exec
DELETE FROM organization_members
WHERE  organization_id = $1 AND member_address = $2;

-- name: GetOrganizationRoles :many
SELECT DISTINCT role FROM organization_members
WHERE organization_id = $1;

-- name: GetAllOrganizationMembers :many
SELECT * FROM organization_members 
WHERE organization_id = $1;

