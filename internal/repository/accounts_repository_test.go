package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestInsertAccount(t *testing.T) {
	t.Run("Completely valid request", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		rows := pgxmock.NewRows([]string{"id", "document_number"}).
			AddRow(int64(1), "12345678900")

		mockDB.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnRows(rows)

		account, err := repo.InsertAccount(ctx, "12345678900")

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, int64(1), account.ID)
		assert.Equal(t, "12345678900", account.DocumentNumber)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Database error during insertion", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		mockDB.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnError(errors.New("database error"))

		account, err := repo.InsertAccount(ctx, "12345678900")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "database error")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Empty document_number should fail", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		mockDB.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("").
			WillReturnError(errors.New("null value in column \"document_number\" violates not-null constraint"))

		account, err := repo.InsertAccount(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "null value in column")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Duplicate document_number should return error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		mockDB.ExpectQuery(`INSERT INTO accounts`).
			WithArgs("12345678900").
			WillReturnError(errors.New("duplicate key value violates unique constraint"))

		account, err := repo.InsertAccount(ctx, "12345678900")

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGetAccountByID(t *testing.T) {
	t.Run("Completely valid request", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(1)

		rows := pgxmock.NewRows([]string{"id", "document_number"}).
			AddRow(accountID, "12345678900")

		mockDB.ExpectQuery(`SELECT id, document_number FROM accounts WHERE id = \$1`).
			WithArgs(accountID).
			WillReturnRows(rows)

		account, err := repo.GetAccountByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, int64(1), account.ID)
		assert.Equal(t, "12345678900", account.DocumentNumber)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Account not found", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(999)

		mockDB.ExpectQuery(`SELECT id, document_number FROM accounts WHERE id = \$1`).
			WithArgs(accountID).
			WillReturnError(pgx.ErrNoRows)

		account, err := repo.GetAccountByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, account)
		assert.Contains(t, err.Error(), "no rows in result set")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
