package repository

import (
	"context"
	"fmt"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/rs/zerolog/log"
)

func NewTransactionsRepository(db PgxPoolIface) TransactionsRepository {
	return &transactionsRepo{db: db}
}

// InsertTransaction inserts a new transaction
func (r *transactionsRepo) InsertTransaction(ctx context.Context, accountID, operationTypeID int64, amount, balance float64) (*Transaction, error) {
	query := `INSERT INTO transactions (account_id, operation_type_id, amount, balance) VALUES ($1, $2, $3, $4) RETURNING id, event_date, balance`
	transaction := &Transaction{}

	err := r.db.QueryRow(ctx, query, accountID, operationTypeID, amount, balance).Scan(
		&transaction.ID,
		&transaction.EventDate,
		&transaction.Balance,
	)
	if err != nil {
		reqID := middleware.GetRequestIDFromContext(ctx)
		log.Error().Str("request_id", reqID).Err(err).Msg("Database error: failed to insert transaction")
		return nil, err
	}

	transaction.AccountID = accountID
	transaction.Amount = amount
	transaction.OperationTypeID = operationTypeID

	return transaction, nil
}

// GetOutstandingTransactionsByAccountID retrieves list of transactions for a given accountID
func (r *transactionsRepo) GetOutstandingTransactionsByAccountID(ctx context.Context, accountID int64) ([]*Transaction, error) {
	var transactions []*Transaction
	query := `SELECT id, amount, balance, event_date 
		FROM transactions 
		WHERE account_id = $1 
		  AND operation_type_id IN (1,2,3) 
		  AND balance < 0 
		ORDER BY event_date`

	rows, err := r.db.Query(ctx,
		query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		txn := &Transaction{}
		if err := rows.Scan(&txn.ID, &txn.Amount, &txn.Balance, &txn.EventDate); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

// UpdateTransactionBalance updates the balance for the provided transactionID
func (r *transactionsRepo) UpdateTransactionBalance(ctx context.Context, transactionID int64, newBalance float64) error {
	query := `UPDATE transactions SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	res, err := r.db.Exec(ctx, query, newBalance, transactionID)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		errMsg := fmt.Errorf("failed to update transaction balance: unexpected number of rows affected: %d for transaction %d", res.RowsAffected(), transactionID)
		log.Error().Err(errMsg).Msg("Database error")
		return errMsg
	}

	log.Info().Msgf("Updated transaction %d with new balance %.2f", transactionID, newBalance)
	return nil
}
