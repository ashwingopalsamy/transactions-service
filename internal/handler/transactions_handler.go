package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/ashwingopalsamy/transactions-service/internal/writer"
)

func NewTransactionHandler(transactionService service.TransactionsService) *TransactionsHandler {
	return &TransactionsHandler{transactionService: transactionService}
}

// CreateTransaction creates new transaction
func (h *TransactionsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionReq

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidRequest,
			ErrInvalidReqBody,
		)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(r.Context(), req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeTransactionErr,
			ErrTitleTrxFailed,
			err.Error(),
		)
		return
	}

	writer.WriteJSON(w, http.StatusCreated, transaction)
	return
}
