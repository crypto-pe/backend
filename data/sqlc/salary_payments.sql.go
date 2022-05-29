// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: salary_payments.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createSalaryPayments = `-- name: CreateSalaryPayments :one

INSERT INTO salary_payments(organization_id, member_address, transaction_hash, amount, token)
VALUES ($1, $2, $3, $4, $5)
RETURNING payment_id, organization_id, member_address, transaction_hash, amount, token, date
`

type CreateSalaryPaymentsParams struct {
	OrganizationID  uuid.UUID `json:"organizationID"`
	MemberAddress   string    `json:"memberAddress"`
	TransactionHash string    `json:"transactionHash"`
	Amount          string    `json:"amount"`
	Token           string    `json:"token"`
}

// payment_id	uuid [uuid_generate_v4()]
// organization_id	uuid
// member_address	bytea
// transaction_hash	character(66)
// amount	integer
// token	bytea
// date	timestamp
func (q *Queries) CreateSalaryPayments(ctx context.Context, arg CreateSalaryPaymentsParams) (SalaryPayments, error) {
	row := q.db.QueryRowContext(ctx, createSalaryPayments,
		arg.OrganizationID,
		arg.MemberAddress,
		arg.TransactionHash,
		arg.Amount,
		arg.Token,
	)
	var i SalaryPayments
	err := row.Scan(
		&i.PaymentID,
		&i.OrganizationID,
		&i.MemberAddress,
		&i.TransactionHash,
		&i.Amount,
		&i.Token,
		&i.Date,
	)
	return i, err
}

const getMemberOverallSalaryHistory = `-- name: GetMemberOverallSalaryHistory :many
SELECT payment_id, organization_id, member_address, transaction_hash, amount, token, date FROM salary_payments
WHERE member_address = $1 ORDER BY date DESC
`

func (q *Queries) GetMemberOverallSalaryHistory(ctx context.Context, memberAddress string) ([]SalaryPayments, error) {
	rows, err := q.db.QueryContext(ctx, getMemberOverallSalaryHistory, memberAddress)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SalaryPayments
	for rows.Next() {
		var i SalaryPayments
		if err := rows.Scan(
			&i.PaymentID,
			&i.OrganizationID,
			&i.MemberAddress,
			&i.TransactionHash,
			&i.Amount,
			&i.Token,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOrgMemberSalaryPaymentsHistory = `-- name: GetOrgMemberSalaryPaymentsHistory :many
SELECT payment_id, organization_id, member_address, transaction_hash, amount, token, date FROM salary_payments
WHERE member_address = $1 AND organization_id = $2
`

type GetOrgMemberSalaryPaymentsHistoryParams struct {
	MemberAddress  string    `json:"memberAddress"`
	OrganizationID uuid.UUID `json:"organizationID"`
}

func (q *Queries) GetOrgMemberSalaryPaymentsHistory(ctx context.Context, arg GetOrgMemberSalaryPaymentsHistoryParams) ([]SalaryPayments, error) {
	rows, err := q.db.QueryContext(ctx, getOrgMemberSalaryPaymentsHistory, arg.MemberAddress, arg.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SalaryPayments
	for rows.Next() {
		var i SalaryPayments
		if err := rows.Scan(
			&i.PaymentID,
			&i.OrganizationID,
			&i.MemberAddress,
			&i.TransactionHash,
			&i.Amount,
			&i.Token,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSalaryPaymentsByTransaction = `-- name: GetSalaryPaymentsByTransaction :many
SELECT payment_id, organization_id, member_address, transaction_hash, amount, token, date FROM salary_payments
WHERE transaction_hash = $1
`

func (q *Queries) GetSalaryPaymentsByTransaction(ctx context.Context, transactionHash string) ([]SalaryPayments, error) {
	rows, err := q.db.QueryContext(ctx, getSalaryPaymentsByTransaction, transactionHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SalaryPayments
	for rows.Next() {
		var i SalaryPayments
		if err := rows.Scan(
			&i.PaymentID,
			&i.OrganizationID,
			&i.MemberAddress,
			&i.TransactionHash,
			&i.Amount,
			&i.Token,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
