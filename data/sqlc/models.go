// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package sqlc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Accounts struct {
	Address   string       `json:"address"`
	Name      string       `json:"name"`
	CreatedAt sql.NullTime `json:"createdAt"`
	Email     interface{}  `json:"email"`
	Admin     sql.NullBool `json:"admin"`
}

type OrganizationMembers struct {
	OrganizationID uuid.UUID      `json:"organizationID"`
	MemberAddress  string         `json:"memberAddress"`
	DateJoined     time.Time      `json:"dateJoined"`
	Role           string         `json:"role"`
	IsAdmin        sql.NullBool   `json:"isAdmin"`
	Salary         sql.NullString `json:"salary"`
}

type Organizations struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	OwnerAddress string    `json:"ownerAddress"`
	Token        string    `json:"token"`
}

type SalaryPayments struct {
	PaymentID       uuid.UUID `json:"paymentID"`
	OrganizationID  uuid.UUID `json:"organizationID"`
	MemberAddress   string    `json:"memberAddress"`
	TransactionHash string    `json:"transactionHash"`
	Amount          string    `json:"amount"`
	Token           string    `json:"token"`
	Date            time.Time `json:"date"`
}
