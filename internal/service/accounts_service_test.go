package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ashwingopalsamy/transactions-service/internal/repository"
	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	t.Run("Valid document number should create account", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		accService := service.NewAccountsService(repo)
		ctx := context.Background()

		rows := pgxmock.NewRows([]string{"id", "document_number"}).AddRow(int64(1), "12345678900")
		mockDB.ExpectQuery(`INSERT INTO accounts`).WithArgs("12345678900").WillReturnRows(rows)

		account, err := accService.CreateAccount(ctx, "12345678900")
		assert.NoError(t, err)
		assert.NotNil(t, account)
	})

	t.Run("Empty document number should fail", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		accService := service.NewAccountsService(repo)
		ctx := context.Background()

		account, err := accService.CreateAccount(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("Valid account ID should return account", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		accService := service.NewAccountsService(repo)
		ctx := context.Background()

		rows := pgxmock.NewRows([]string{"id", "document_number"}).AddRow(int64(1), "12345678900")
		mockDB.ExpectQuery(`SELECT id, document_number FROM accounts WHERE id = \$1`).WithArgs(int64(1)).WillReturnRows(rows)

		account, err := accService.GetAccount(ctx, 1)
		assert.NoError(t, err)
		assert.NotNil(t, account)
	})

	t.Run("Non-existing account ID should return error", func(t *testing.T) {
		mockDB, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mockDB.Close()

		repo := repository.NewAccountsRepository(mockDB)
		accService := service.NewAccountsService(repo)
		ctx := context.Background()

		mockDB.ExpectQuery(`SELECT id, document_number FROM accounts WHERE id = \$1`).WithArgs(int64(999)).WillReturnError(errors.New("no rows in result set"))

		account, err := accService.GetAccount(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, account)
	})
}
