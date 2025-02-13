package repository_test

import (
	"context"
	"errors"
	"fmt"
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
		balance := amount

		rows := pgxmock.NewRows([]string{"id", "event_date", "balance"}).
			AddRow(int64(1), time.Now(), balance)

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, operationTypeID, amount, balance).
			WillReturnRows(rows)

		transaction, err := repo.InsertTransaction(ctx, accountID, operationTypeID, amount, balance)

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
		balance := amount

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, operationTypeID, amount, balance).
			WillReturnError(errors.New("database error"))

		transaction, err := repo.InsertTransaction(ctx, accountID, operationTypeID, amount, balance)

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
		balance := amount

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(invalidAccountID, operationTypeID, amount, balance).
			WillReturnError(errors.New("violates foreign key constraint \"transactions_account_id_fkey\""))

		transaction, err := repo.InsertTransaction(ctx, invalidAccountID, operationTypeID, amount, balance)

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
		balance := amount

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(accountID, invalidOperationTypeID, amount, balance).
			WillReturnError(errors.New("violates foreign key constraint \"transactions_operation_type_id_fkey\""))

		transaction, err := repo.InsertTransaction(ctx, accountID, invalidOperationTypeID, amount, balance)

		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Contains(t, err.Error(), "violates foreign key constraint \"transactions_operation_type_id_fkey\"")

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestGetOutstandingTransactionsByAccountID(t *testing.T) {
	t.Run("Successful retrieval of transactions by accountID", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		accountID := int64(1)

		now := time.Now()
		later := now.Add(10 * time.Second)
		rows := pgxmock.NewRows([]string{"id", "amount", "balance", "event_date"}).AddRow(int64(1), 100.00, 100.00, now).
			AddRow(int64(2), 100.00, 100.00, later)

		mockDB.ExpectQuery(`SELECT id, amount, balance, event_date FROM transactions WHERE account_id = \$1`).WithArgs(accountID).WillReturnRows(rows)

		txns, err := repo.GetOutstandingTransactionsByAccountID(ctx, accountID)
		assert.NoError(t, err)
		assert.Len(t, txns, 2)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

}

func TestUpdateTransactionBalance(t *testing.T) {
	t.Run("Successful update returns nil error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		transactionID := int64(123)
		newBalance := 500.25

		mockDB.ExpectExec(`UPDATE transactions SET balance = \$1, updated_at = CURRENT_TIMESTAMP WHERE id = \$2`).
			WithArgs(newBalance, transactionID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = repo.UpdateTransactionBalance(ctx, transactionID, newBalance)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Unexpected rows affected returns error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewTransactionsRepository(mockDB)
		ctx := context.Background()

		transactionID := int64(789)
		newBalance := 1000.00

		mockDB.ExpectExec(`UPDATE transactions SET balance = \$1, updated_at = CURRENT_TIMESTAMP WHERE id = \$2`).
			WithArgs(newBalance, transactionID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err = repo.UpdateTransactionBalance(ctx, transactionID, newBalance)
		assert.Error(t, err)
		expectedErrMsg := fmt.Sprintf("failed to update transaction balance: unexpected number of rows affected: %d for transaction %d", 0, transactionID)
		assert.EqualError(t, err, expectedErrMsg)

		assert.NoError(t, mockDB.ExpectationsWereMet())
	})
}
