package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	CorrelationIDKey string = "Correlation-ID"
)

func (m *Middleware) LogContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CorrelationIDKey, uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
