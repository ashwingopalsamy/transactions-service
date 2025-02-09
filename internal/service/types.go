package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
)

type AccountsService interface {
	CreateAccount(ctx context.Context, documentNumber string) (*repository.Account, error)
	GetAccount(ctx context.Context, accountID int64) (*repository.Account, error)
}

type TransactionsService interface {
	CreateTransaction(ctx context.Context, accountID, operationTypeID int64, amount float64) (*repository.Transaction, error)
}

type accountsService struct {
	accRepo repository.AccountsRepository
}

type transactionsService struct {
	trxRepo repository.TransactionsRepository
	accRepo repository.AccountsRepository
}

// Account-related errors
var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrAccountAlreadyExists  = errors.New("document_number already exists")
	ErrInvalidDocumentNumber = errors.New("document_number cannot be empty")
	ErrFailedToFetchAccount  = errors.New("failed to fetch account")
)

// Transaction-related errors
var (
	ErrInvalidAccountID     = errors.New("invalid account_id: account does not exist")
	ErrInvalidOperationType = errors.New("invalid operation_type_id: operation type does not exist")
	ErrInvalidAmount        = errors.New("invalid amount: amount must not be zero")
	ErrNegativeAmount       = errors.New("invalid amount: amount must not be negative")
	ErrTransactionFailed    = errors.New("failed to insert transaction")
)

// determinePgxError maps pgx constraint violations to known errors.
func determinePgxError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	if strings.Contains(errMsg, "violates foreign key constraint") {
		if strings.Contains(errMsg, "transactions_account_id_fkey") {
			return ErrInvalidAccountID
		}
		if strings.Contains(errMsg, "transactions_operation_type_id_fkey") {
			return ErrInvalidOperationType
		}
	}

	if strings.Contains(errMsg, "unique constraint") {
		return ErrAccountAlreadyExists
	}

	return err
}
