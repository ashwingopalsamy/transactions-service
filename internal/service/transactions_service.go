package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/jackc/pgx/v5"
)

func NewTransactionsService(trxRepo repository.TransactionsRepository, accRepo repository.AccountsRepository) TransactionsService {
	return &transactionsService{trxRepo: trxRepo, accRepo: accRepo}
}

// CreateTransaction validates and creates a transaction
func (s *transactionsService) CreateTransaction(ctx context.Context, accountID, operationTypeID int64, amount float64) (*repository.Transaction, error) {
	// Check if the account exists
	_, err := s.accRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidAccountID
		}
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	if amount == 0 {
		return nil, ErrInvalidAmount
	}

	// Ensure amount is appropriately signed
	amount, err = EnforceAmountSign(operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	// Format the amount to have exactly two decimal places.
	amount = FormatAmount(amount)

	// Insert transaction record
	transaction, err := s.trxRepo.InsertTransaction(ctx, accountID, operationTypeID, amount)
	if err != nil {
		return nil, determinePgxError(err)
	}

	return transaction, nil
}

// EnforceAmountSign ensures that certain transaction types have positive/negative amounts
func EnforceAmountSign(operationTypeID int64, amount float64) (float64, error) {
	switch operationTypeID {
	case 1, 2, 3: // Purchases and withdrawals → Negative amount
		if amount > 0 {
			amount = -amount
		}
	case 4: // Credit Voucher → Positive amount
		if amount < 0 {
			amount = -amount
		}
	default:
		return 0, ErrInvalidOperationType
	}
	return amount, nil
}

// FormatAmount ensures the float has exactly two decimals. .
func FormatAmount(amount float64) float64 {
	formatted := fmt.Sprintf("%.2f", amount)

	var result float64
	fmt.Sscanf(formatted, "%f", &result)
	return result
}
