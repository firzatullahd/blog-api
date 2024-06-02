package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/firzatullahd/blog-api/internal/model/response"
	log "github.com/sirupsen/logrus"
)

type Middleware struct {
	JWTSecretKey string
}

func NewMiddleware(jwtSecretKey string) *Middleware {
	return &Middleware{
		JWTSecretKey: jwtSecretKey,
	}
}

func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v %v", err, string(debug.Stack()))
				response.SetHTTPResponse(w, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		next.ServeHTTP(w, req)
	})

}
