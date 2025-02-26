package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AccountsRepository interface {
	InsertAccount(ctx context.Context, documentNumber string) (*Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*Account, error)
}

type TransactionsRepository interface {
	InsertTransaction(ctx context.Context, accountID, operationTypeID int64, amount, balance float64) (*Transaction, error)

	GetOutstandingTransactionsByAccountID(ctx context.Context, accountID int64) ([]*Transaction, error)
	UpdateTransactionBalance(ctx context.Context, transactionID int64, amount float64) error
}

type PgxPoolIface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type accountsRepo struct {
	db PgxPoolIface
}

type transactionsRepo struct {
	db PgxPoolIface
}

type Account struct {
	ID             int64     `json:"id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"-"`
}

// Transaction
// Ideally, we should not be treating currencies ('Amount') in float64
// ref: https://stackoverflow.com/questions/3730019/why-not-use-double-or-float-to-represent-currency
// To satisfy the tech-case requirements, the amount is defined in float64
type Transaction struct {
	ID              int64     `json:"id"`
	AccountID       int64     `json:"-"`
	OperationTypeID int64     `json:"-"`
	Amount          float64   `json:"-"`
	Balance         float64   `json:"-"`
	EventDate       time.Time `json:"event_date"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
