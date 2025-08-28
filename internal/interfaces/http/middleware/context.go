package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	// RequestIDKey is the key used to store the request ID in the context
	RequestIDKey contextKey = "requestID"
)

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return uuid.New().String()
}

// RequestIDMiddleware adds a request ID to the context and response headers
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate or get request ID from header
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// In a real app, you might want to generate a UUID here
			requestID = uuid.NewString()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetPaginationParams extracts and validates pagination parameters from the request
func GetPaginationParams(r *http.Request) (int64, int64, error) {
	page := int64(1)
	limit := int64(10)

	// Parse page parameter
	if p := r.URL.Query().Get("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil {
			return 0, 0, fmt.Errorf("invalid page parameter")
		}
		if page < 1 {
			page = 1
		}
	}

	// Parse limit parameter
	if l := r.URL.Query().Get("limit"); l != "" {
		if _, err := fmt.Sscanf(l, "%d", &limit); err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter")
		}
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 10000000
		}
	}

	return page, limit, nil
}
