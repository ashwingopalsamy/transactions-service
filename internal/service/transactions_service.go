package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
)

func NewTransactionsService(trxRepo repository.TransactionsRepository, accRepo repository.AccountsRepository) TransactionsService {
	return &transactionsService{trxRepo: trxRepo, accRepo: accRepo}
}

// CreateTransaction validates and creates a transaction
func (s *transactionsService) CreateTransaction(ctx context.Context, accountID, operationTypeID int64, amount float64) (*repository.Transaction, error) {
	// Check if the account exists
	_, err := s.accRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, errors.New("invalid account_id: account does not exist")
	}

	if amount == 0 {
		return nil, errors.New("invalid amount: amount must not be zero")
	}

	// Ensure amount is appropriately signed
	amount, err = s.enforceAmountSign(operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	// Format the amount to have exactly two decimal places.
	amount = formatAmount(amount)

	// Insert transaction record
	transaction, err := s.trxRepo.InsertTransaction(ctx, accountID, operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// enforceAmountSign ensures that certain transaction types have positive/negative amounts
func (s *transactionsService) enforceAmountSign(operationTypeID int64, amount float64) (float64, error) {
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
		return 0, errors.New("invalid operation_type_id: operation type does not exist")
	}
	return amount, nil
}

// formatAmount ensures the float has exactly two decimals. .
func formatAmount(amount float64) float64 {
	formatted := fmt.Sprintf("%.2f", amount)

	var result float64
	fmt.Sscanf(formatted, "%f", &result)
	return result
}
