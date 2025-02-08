package repository

import (
	"context"
	"fmt"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/rs/zerolog/log"
)

func NewAccountsRepository(db PgxPoolIface) AccountsRepository {
	return &accountsRepo{db: db}
}

// InsertAccount inserts a new account
func (r *accountsRepo) InsertAccount(ctx context.Context, documentNumber string) (*Account, error) {
	query := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING id, document_number`
	account := &Account{}

	err := r.db.QueryRow(ctx, query, documentNumber).Scan(&account.ID, &account.DocumentNumber)
	if err != nil {
		reqID := middleware.GetRequestIDFromContext(ctx)
		log.Error().Str("request_id", reqID).Err(err).Msg("Database error: failed to insert account")
		return nil, fmt.Errorf("failed to insert account: %w", err)
	}
	return account, nil
}

// GetAccountByID retrieves an account by accountID
func (r *accountsRepo) GetAccountByID(ctx context.Context, accountID int64) (*Account, error) {
	query := `SELECT id, document_number FROM accounts WHERE id = $1`
	account := &Account{}

	err := r.db.QueryRow(ctx, query, accountID).Scan(&account.ID, &account.DocumentNumber)
	if err != nil {
		reqID := middleware.GetRequestIDFromContext(ctx)
		log.Error().Str("request_id", reqID).Err(err).Msg("Database error: failed to retrieve account")
		return nil, err
	}
	return account, nil
}
