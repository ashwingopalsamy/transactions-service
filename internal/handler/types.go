package handler

import (
	"github.com/ashwingopalsamy/transactions-service/internal/service"
)

// Error Codes, Titles and Error Details
const (
	ErrCodeInvalidRequest = "invalid_request"
	ErrCodeConflictErr    = "conflict_error"
	ErrCodeTransactionErr = "transaction_error"

	ErrTitleAccNotFound    = "Account Not Found"
	ErrTitleConflict       = "Conflict"
	ErrTitleInvalidAccID   = "Invalid Account ID"
	ErrTitleInvalidRequest = "Invalid Request"
	ErrTitleTrxFailed      = "Transaction Failed"

	ErrInvalidReqBody = "invalid request body"
)

type AccountsHandler struct {
	accountService service.AccountsService
}

type TransactionsHandler struct {
	transactionService service.TransactionsService
}

type CreateAccountReq struct {
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionReq struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}
