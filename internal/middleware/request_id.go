package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type key int

var requestIDKey = key(1)

// SetRequestIDToContext creates a UUIDv7-based requestID
func SetRequestIDToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := uuid.NewV7()
		reqIDStr := reqID.String()

		// Set the header
		w.Header().Set(middleware.RequestIDHeader, reqIDStr)

		// Set requestID to the context.
		ctx := context.WithValue(r.Context(), requestIDKey, reqIDStr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestIDFromContext retrieves the requestID from context
func GetRequestIDFromContext(ctx context.Context) string {
	reqID, _ := ctx.Value(requestIDKey).(string)
	return reqID
}
