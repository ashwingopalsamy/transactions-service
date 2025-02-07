package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	t.Run("Valid transaction should succeed", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		trxRepo := repository.NewTransactionsRepository(mockDB)
		accRepo := repository.NewAccountsRepository(mockDB)
		trxService := service.NewTransactionsService(trxRepo, accRepo)
		ctx := context.Background()

		mockDB.ExpectQuery(`SELECT id, document_number FROM accounts WHERE id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "document_number"}).
				AddRow(int64(1), "12345678900"))

		mockDB.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(int64(1), int64(4), float64(100.00)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "event_date"}).
				AddRow(int64(1), time.Now()))

		transaction, err := trxService.CreateTransaction(ctx, int64(1), 4, 100.00)
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
	})

	t.Run("Zero amount should fail", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		trxRepo := repository.NewTransactionsRepository(mockDB)
		accRepo := repository.NewAccountsRepository(mockDB)
		trxService := service.NewTransactionsService(trxRepo, accRepo)
		ctx := context.Background()

		transaction, err := trxService.CreateTransaction(ctx, 1, 4, 0)
		assert.Error(t, err)
		assert.Nil(t, transaction)
	})

	t.Run("Invalid operation type should fail", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		trxRepo := repository.NewTransactionsRepository(mockDB)
		accRepo := repository.NewAccountsRepository(mockDB)
		trxService := service.NewTransactionsService(trxRepo, accRepo)
		ctx := context.Background()

		transaction, err := trxService.CreateTransaction(ctx, 1, 99, 100.00)
		assert.Error(t, err)
		assert.Nil(t, transaction)
	})
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name           string
		inputAmount    float64
		expectedAmount float64
	}{
		{"Round down", 123.456, 123.46},
		{"Round up", 99.999, 100.00},
		{"Exact two decimals", 50.50, 50.50},
		{"Zero value", 0.00, 0.00},
		{"Negative value rounding", -25.555, -25.55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount := service.FormatAmount(tt.inputAmount)
			assert.Equal(t, tt.expectedAmount, amount)
		})
	}
}

func TestEnforceAmountSign(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int64
		amount          float64
		expectedAmount  float64
		expectedError   error
	}{
		{"Normal Purchase - Positive to Negative", 1, 100.00, -100.00, nil},
		{"Purchase with Installments - Positive to Negative", 2, 50.00, -50.00, nil},
		{"Withdrawal - Positive to Negative", 3, 25.50, -25.50, nil},
		{"Credit Voucher - Negative to Positive", 4, -75.25, 75.25, nil},
		{"Normal Purchase - Negative Remains Negative", 1, -200.00, -200.00, nil},
		{"Credit Voucher - Positive Remains Positive", 4, 150.00, 150.00, nil},
		{"Invalid Operation Type", 99, 100.00, 0, errors.New("invalid operation_type_id: operation type does not exist")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := service.EnforceAmountSign(tt.operationTypeID, tt.amount)
			assert.Equal(t, tt.expectedAmount, amount)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
