-- payment_id	uuid [uuid_generate_v4()]	
-- organization_id	uuid	
-- member_address	bytea	
-- transaction_hash	character(66)	
-- amount	integer	
-- token	bytea	
-- date	timestamp

-- name: CreateSalaryPayment :one
INSERT INTO salary_payments(organization_id, member_address, transaction_hash, amount, token, date)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSalaryPayment :one
SELECT * FROM salary_payments
WHERE transaction_hash = $1;

--name: GetAllMemberSalaryPayments :many
SELECT * FROM salary_payments
WHERE member_address = $1;
