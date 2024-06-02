package route

import (
	"net/http"

	"github.com/firzatullahd/blog-api/internal/config"
	"github.com/firzatullahd/blog-api/internal/delivery/http/handler"
	"github.com/firzatullahd/blog-api/internal/delivery/http/middleware"
	"github.com/firzatullahd/blog-api/internal/model/response"
	log "github.com/sirupsen/logrus"
)

func Serve(conf *config.Config, h *handler.Handler) {

	m := middleware.NewMiddleware(conf.JWTSecretKey)
	mux := http.NewServeMux()

	mux.Handle("GET /health", m.LogContext(http.HandlerFunc(Health)))
	mux.Handle("POST /api/users/register", m.LogContext(http.HandlerFunc(h.Register)))
	mux.Handle("POST /api/users/login", m.LogContext(http.HandlerFunc(h.Login)))
	mux.Handle("POST /api/users/admin", m.LogContext(http.HandlerFunc(h.GrantAdmin)))

	mux.Handle("POST /api/posts", m.LogContext(m.Auth(http.HandlerFunc(h.CreatePost))))
	mux.Handle("PUT /api/posts/{id}", m.LogContext(m.Auth(http.HandlerFunc(h.UpdatePost))))
	mux.Handle("DELETE /api/posts/{id}", m.LogContext(m.Auth(http.HandlerFunc(h.DeletePost))))
	mux.Handle("GET /api/posts/{id}", m.LogContext(m.Auth(http.HandlerFunc(h.GetPost))))

	log.Println("Running on port " + conf.Port)
	log.Fatal(http.ListenAndServe(conf.Port, mux))

}

func Health(w http.ResponseWriter, r *http.Request) {
	response.SetHTTPResponse(w, http.StatusOK, "success", nil)
}
