package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func NewTransactionsRepository(db PgxPoolIface) TransactionsRepository {
	return &transactionsRepo{db: db}
}

// InsertTransaction inserts a new transaction
func (r *transactionsRepo) InsertTransaction(ctx context.Context, accountID, operationTypeID int64, amount float64) (*Transaction, error) {
	query := `INSERT INTO transactions (account_id, operation_type_id, amount) VALUES ($1, $2, $3) RETURNING id, event_date`
	transaction := &Transaction{}

	err := r.db.QueryRow(ctx, query, accountID, operationTypeID, amount).Scan(
		&transaction.ID,
		&transaction.EventDate,
	)
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			if strings.Contains(err.Error(), "transactions_account_id_fkey") {
				return nil, errors.New("invalid account_id: account does not exist")
			}
			if strings.Contains(err.Error(), "transactions_operation_type_id_fkey") {
				return nil, errors.New("invalid operation_type_id: operation type does not exist")
			}
			return nil, errors.New("invalid foreign key reference")
		}
		return nil, fmt.Errorf("failed to insert transaction: %w", err)
	}
	return transaction, nil
}
