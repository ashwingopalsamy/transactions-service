package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestInsertTransaction(t *testing.T) {
	t.Run("Completely valid transaction", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(1)
		operationTypeID := int64(4)
		amount := 100.50

		rows := pgxmock.NewRows([]string{"id", "event_date"}).
			AddRow(int64(1), time.Now())

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, operationTypeID, amount).
			WillReturnRows(rows)

		transaction, err := repo.InsertTransaction(ctx, accountID, operationTypeID, amount)

		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, int64(1), transaction.ID)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Database error during insertion", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(1)
		operationTypeID := int64(4)
		amount := 100.50

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, operationTypeID, amount).
			WillReturnError(errors.New("database error"))

		transaction, err := repo.InsertTransaction(ctx, accountID, operationTypeID, amount)

		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Contains(t, err.Error(), "database error")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Invalid account_id should return error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		invalidAccountID := int64(999)
		operationTypeID := int64(4)
		amount := 100.50

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(invalidAccountID, operationTypeID, amount).
			WillReturnError(errors.New("violates foreign key constraint \"transactions_account_id_fkey\""))

		transaction, err := repo.InsertTransaction(ctx, invalidAccountID, operationTypeID, amount)

		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Contains(t, err.Error(), "violates foreign key constraint \"transactions_account_id_fkey\"")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Invalid operation_type_id should return error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(1)
		invalidOperationTypeID := int64(99)
		amount := 100.50

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, invalidOperationTypeID, amount).
			WillReturnError(errors.New("violates foreign key constraint \"transactions_operation_type_id_fkey\""))

		transaction, err := repo.InsertTransaction(ctx, accountID, invalidOperationTypeID, amount)

		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Contains(t, err.Error(), "violates foreign key constraint \"transactions_operation_type_id_fkey\"")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
