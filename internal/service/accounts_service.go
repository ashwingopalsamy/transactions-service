package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/jackc/pgx/v5"
)

func NewAccountsService(accRepo repository.AccountsRepository) AccountsService {
	return &accountsService{accRepo: accRepo}
}

// CreateAccount creates a new account
func (s *accountsService) CreateAccount(ctx context.Context, documentNumber string) (*repository.Account, error) {
	if strings.TrimSpace(documentNumber) == "" {
		return nil, errors.New("document_number cannot be empty")
	}

	account, err := s.accRepo.InsertAccount(ctx, documentNumber)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccount retrieves an account by accountID
func (s *accountsService) GetAccount(ctx context.Context, accountID int64) (*repository.Account, error) {
	account, err := s.accRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, errors.New("failed to fetch account")
	}
	return account, nil
}
