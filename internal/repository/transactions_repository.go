package repository

import (
	"context"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/rs/zerolog/log"
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
		reqID := middleware.GetRequestIDFromContext(ctx)
		log.Error().Str("request_id", reqID).Err(err).Msg("Database error: failed to insert transaction")
		return nil, err
	}
	return transaction, nil
}
