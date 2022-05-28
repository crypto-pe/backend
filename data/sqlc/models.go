// Code generated by sqlc. DO NOT EDIT.

package sqlc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Accounts struct {
	Address   []byte       `json:"address"`
	Name      string       `json:"name"`
	CreatedAt sql.NullTime `json:"createdAt"`
	Email     interface{}  `json:"email"`
	Admin     sql.NullBool `json:"admin"`
}

type OrganizationMembers struct {
	OrganizationID uuid.UUID     `json:"organizationID"`
	MemberAddress  []byte        `json:"memberAddress"`
	DateJoined     time.Time     `json:"dateJoined"`
	Role           string        `json:"role"`
	IsAdmin        sql.NullBool  `json:"isAdmin"`
	Salary         sql.NullInt32 `json:"salary"`
}

type Organizations struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	OwnerAddress []byte    `json:"ownerAddress"`
	Token        []byte    `json:"token"`
}

type SalaryPayments struct {
	PaymentID       uuid.UUID `json:"paymentID"`
	OrganizationID  uuid.UUID `json:"organizationID"`
	MemberAddress   []byte    `json:"memberAddress"`
	TransactionHash string    `json:"transactionHash"`
	Amount          int32     `json:"amount"`
	Token           []byte    `json:"token"`
	Date            time.Time `json:"date"`
}