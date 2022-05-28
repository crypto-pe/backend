-- organization_id	uuid	
-- member_address	bytea	
-- date_joined	timestamp	
-- role	character varying(200)	
-- is_admin	boolean NULL [false]	
-- salary	integer NULL [0]	

-- name: CreateOrganizationMember :one
INSERT INTO organization_members(
  organization_id, member_address, date_joined, role, is_admin, salary
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetOrgamizationMember :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND member_address = $2;

-- name: UpdateOrganizationMember :one
UPDATE organization_members
SET
  date_joined = $3,
  role = $4, 
  is_admin = $5,
  salary = $6
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

