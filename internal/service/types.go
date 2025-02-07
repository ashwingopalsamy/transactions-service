package service

import (
	"context"

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
