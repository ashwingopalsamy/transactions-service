package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func NewTransactionsService(trxRepo repository.TransactionsRepository, accRepo repository.AccountsRepository) TransactionsService {
	return &transactionsService{trxRepo: trxRepo, accRepo: accRepo}
}

// CreateTransaction validates and creates a transaction
func (s *transactionsService) CreateTransaction(ctx context.Context, accountID, operationTypeID int64, amount float64) (*repository.Transaction, error) {
	// Check if the account exists
	if _, err := s.accRepo.GetAccountByID(ctx, accountID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidAccountID
		}
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}

	// Validate amount: must be strictly positive.
	if amount <= 0 {
		if amount == 0 {
			return nil, ErrInvalidAmount
		}
		return nil, ErrNegativeAmount
	}

	// Ensure amount is appropriately signed
	amount, err := EnforceAmountSign(operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	// Format the amount to have exactly two decimal places.
	amount = FormatAmount(amount)

	// Set balance as amount (initially)
	balance := amount

	// Insert transaction record
	transaction, err := s.trxRepo.InsertTransaction(ctx, accountID, operationTypeID, amount, balance)
	if err != nil {
		return nil, determinePgxError(err)
	}

	// Process Payment Discharge
	// when a credit transaction is found
	if operationTypeID == 4 {
		err := s.processPaymentDischarge(ctx, transaction)
		if err != nil {
			return nil, fmt.Errorf("payment discharge error: %w", err)
		}
	}

	return transaction, nil
}

// processPaymentDischarge applies a payment transaction against outstanding purchase/withdrawal transactions.
func (s *transactionsService) processPaymentDischarge(ctx context.Context, creditTxn *repository.Transaction) error {
	// PLAN
	// 1. Initialize the credit amount (op.type 4).
	// 2. Define the total dischargeable amount as the credit.
	// 3. Retrieve outstanding transactions for the account.
	// 4. Iterate over the transactions to apply payment discharge.
	// 5. Update each outstanding transaction’s balance as discharge is applied.
	// 6. Continue until the credit is fully allocated.
	// 7. Update the op.type 4 transaction with its new balance.

	creditedAmount := creditTxn.Amount
	log.Info().Msgf("Starting Payment Discharge for Txn: %d: creditedAmount = %.2f", creditTxn.ID, creditedAmount)

	// Fetch all outstanding balances for the account.
	outstandingTxns, err := s.trxRepo.GetOutstandingTransactionsByAccountID(ctx, creditTxn.AccountID)
	if err != nil {
		fmt.Println("101")
		fmt.Println(err)
		return err
	}

	var totalDischarge float64
	for _, outstandingTxn := range outstandingTxns {
		if creditedAmount <= 0 {
			log.Info().Msg("no more credited amount left to process discharge")
			break
		}

		// Calculate the total absolute outstanding amount per txn
		outstandingAbsValue := -outstandingTxn.Balance
		log.Info().Msgf("processing discharge for outstanding txn: %d, current outstanding balance: %.2f, current outstanding amount: %.2f", outstandingTxn.ID, outstandingTxn.Balance, outstandingAbsValue)

		// Calculate the dischargeable amount for this transaction
		dischargeableAmount := creditedAmount
		if creditedAmount > outstandingAbsValue {
			dischargeableAmount = outstandingAbsValue
		}

		// Calculate the new balance for the outstanding transaction
		newBalance := outstandingTxn.Balance + dischargeableAmount
		if err := s.trxRepo.UpdateTransactionBalance(ctx, outstandingTxn.ID, newBalance); err != nil {
			fmt.Println("126")
			fmt.Println(err)
			return err
		}

		// After processing, update the total discharged amount,
		// adjust outstanding transaction balances, and finalize the credit transaction balance.
		creditedAmount -= dischargeableAmount
		totalDischarge += dischargeableAmount
	}

	// At-last, update the credit transaction's latest balance value
	// after the payment discharge process is completed
	newBalanceForCreditTxn := creditTxn.Amount - totalDischarge
	if err := s.trxRepo.UpdateTransactionBalance(ctx, creditTxn.ID, newBalanceForCreditTxn); err != nil {
		fmt.Println("137")
		return err
	}

	creditTxn.Balance = newBalanceForCreditTxn
	log.Info().Msgf("Finished payment discharge for txn %d; total discharged = %.2f, final payment balance = %.2f",
		creditTxn.ID, totalDischarge, newBalanceForCreditTxn)
	return nil
}

// EnforceAmountSign ensures that certain transaction types have positive/negative amounts
func EnforceAmountSign(operationTypeID int64, amount float64) (float64, error) {
	switch operationTypeID {
	case 1, 2, 3: // Purchases and withdrawals → Negative amount
		return -math.Abs(amount), nil
	case 4: // // Credit Voucher → Positive amount
		return math.Abs(amount), nil
	default:
		return 0, ErrInvalidOperationType
	}
}

// FormatAmount ensures the float has exactly two decimals. .
func FormatAmount(amount float64) float64 {
	formatted := fmt.Sprintf("%.2f", amount)

	var result float64
	fmt.Sscanf(formatted, "%f", &result)
	return result
}
