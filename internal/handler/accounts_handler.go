package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/ashwingopalsamy/transactions-service/internal/writer"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func NewAccountsHandler(accountService service.AccountsService) *AccountsHandler {
	return &AccountsHandler{accountService: accountService}
}

// CreateAccount handles account creation requests
func (h *AccountsHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestIDFromContext(r.Context())

	var req CreateAccountReq

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Error().Str("request_id", reqID).Err(err).Msg("error decoding create account request")
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
		log.Error().Str("request_id", reqID).Err(fmt.Errorf("invalid create account request")).Msg("invalid request body")
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidRequest,
			ErrInvalidReqBody,
		)
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), req.DocumentNumber)
	if err != nil {
		log.Error().Str("request_id", reqID).Err(err).Msg("failed to create account")
		if strings.Contains(err.Error(), "null value in column") {
			writer.WriteError(
				w, r.Context(),
				http.StatusBadRequest,
				ErrCodeInvalidRequest,
				ErrTitleInvalidRequest,
				err.Error(),
			)
			return
		}
		if strings.Contains(err.Error(), "unique constraint") {
			writer.WriteError(
				w, r.Context(),
				http.StatusConflict,
				ErrCodeConflictErr,
				ErrTitleConflict,
				err.Error(),
			)
			return
		}

		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidRequest,
			err.Error(),
		)
		return
	}

	log.Info().Str("request_id", reqID).Int64("id", account.ID).Msg("account creation successful")
	writer.WriteJSON(w, http.StatusCreated, account)
	return
}

// GetAccount handles retrieving an account by ID
func (h *AccountsHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetRequestIDFromContext(r.Context())

	if r.ContentLength > 0 {
		log.Error().Str("request_id", reqID).Err(fmt.Errorf("invalid request")).Msg("invalid request body")
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidRequest,
			ErrInvalidReqBody,
		)
		return
	}

	accountID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error().Str("request_id", reqID).Err(fmt.Errorf("invalid request")).Msg("invalid request param")
		writer.WriteError(
			w, r.Context(),
			http.StatusBadRequest,
			ErrCodeInvalidRequest,
			ErrTitleInvalidAccID,
			err.Error(),
		)
		return
	}

	account, err := h.accountService.GetAccount(r.Context(), accountID)
	if err != nil {
		log.Error().Str("request_id", reqID).Err(err).Msg("failed to get account")
		writer.WriteError(
			w, r.Context(),
			http.StatusNotFound,
			ErrCodeInvalidRequest,
			ErrTitleAccNotFound,
			err.Error(),
		)
		return
	}

	log.Info().Str("request_id", reqID).Int64("id", account.ID).Msg("account retrieval successful")
	writer.WriteJSON(w, http.StatusOK, account)
	return
}
