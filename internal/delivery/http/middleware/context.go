package middleware

import (
	"context"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
	"net/http"

	"github.com/google/uuid"
)

func (m *Middleware) LogContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), logger.CorrelationIDKey, uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
