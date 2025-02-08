package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/ashwingopalsamy/transactions-service/internal/writer"
	"github.com/rs/zerolog/log"
)

func NewTransactionHandler(transactionService service.TransactionsService) *TransactionsHandler {
	return &TransactionsHandler{transactionService: transactionService}
}

// CreateTransaction creates new transaction
func (h *TransactionsHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestIDFromContext(r.Context())

	var req CreateTransactionReq

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Error().Str("request_id", reqID).Err(err).Msg("error decoding create transaction request")
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidRequest,
			ErrInvalidReqBody,
		)
		return
	}

	if r.ContentLength <= 0 {
		log.Error().Str("request_id", reqID).Err(fmt.Errorf("invalid create transaction request")).Msg("invalid request body")
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
		log.Error().Str("request_id", reqID).Err(err).Msg("failed to create transaction")
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeTransactionErr,
			ErrTitleTrxFailed,
			err.Error(),
		)
		return
	}

	log.Info().Str("request_id", reqID).Int64("id", transaction.ID).Msg("transaction successful")
	writer.WriteJSON(w, http.StatusCreated, transaction)
	return
}
