package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ashwingopalsamy/transactions-service/internal/service"
	"github.com/ashwingopalsamy/transactions-service/internal/writer"
	"github.com/go-chi/chi/v5"
)

func NewAccountsHandler(accountService service.AccountsService) *AccountsHandler {
	return &AccountsHandler{accountService: accountService}
}

// CreateAccount handles account creation requests
func (h *AccountsHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountReq

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

	if r.ContentLength <= 0 {
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

	writer.WriteJSON(w, http.StatusCreated, account)
	return
}

// GetAccount handles retrieving an account by ID
func (h *AccountsHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
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
		writer.WriteError(
			w, r.Context(),
			http.StatusNotFound,
			ErrCodeInvalidRequest,
			ErrTitleAccNotFound,
			err.Error(),
		)
		return
	}

	writer.WriteJSON(w, http.StatusOK, account)
	return
}
