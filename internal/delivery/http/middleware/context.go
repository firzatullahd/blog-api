package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	CorrelationIDKey string = "Correlation-ID"
)

func (m *Middleware) LogContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CorrelationIDKey, uuid.New().String())
		fmt.Println("INVOKE LOGCONTEXT")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
