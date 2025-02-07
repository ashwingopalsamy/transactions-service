package writer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ashwingopalsamy/transactions-service/internal/middleware"
)

// ErrorResponse structures our errors
type ErrorResponse struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// WriteError writes error responses in JSON
func WriteError(w http.ResponseWriter, ctx context.Context, status int, code, title, detail string) {
	reqID := middleware.GetRequestIDFromContext(ctx)
	errorResp := ErrorResponse{
		ID:     reqID,
		Code:   code,
		Status: status,
		Title:  title,
		Detail: detail,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
	}
}

// WriteJSON writes successful responses in JSON with response body
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode success response", http.StatusInternalServerError)
	}
}
