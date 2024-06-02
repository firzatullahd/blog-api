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

	//e := echo.New()
	m := middleware.NewMiddleware(conf.JWTSecretKey)
	//e.Pre(echoMiddleware.RemoveTrailingSlash())
	//e.Use(echoMiddleware.Logger(), echoMiddleware.Recover())
	//e.Use(m.LogContext())

	mux := http.NewServeMux()

	mux.Handle("/", m.LogContext(http.HandlerFunc(hello)))
	mux.Handle("/v1/user/register", m.LogContext(http.HandlerFunc(h.Register)))

	log.Fatal(http.ListenAndServe(conf.Port, mux))

	//userApi := e.Group("/v1/user")
	//userApi.POST("/register", h.Register)
	//userApi.POST("/login", h.Login)
	//
	//catApi := e.Group("/v1/cat", m.Auth())
	//catApi.POST("/", h.CreateCat)
	//catApi.GET("/", h.FindCat)
	//catApi.PUT("/:id", h.UpdateCat)
	//catApi.DELETE("/:id", h.DeleteCat)
	//
	//matchApi := e.Group("/v1/match", m.Auth())
	//matchApi.POST("/", h.CreateMatch)
	//matchApi.GET("/", h.FindMatch)
	//matchApi.POST("/approve", h.ApproveMatch)
	//matchApi.POST("/reject", h.RejectMatch)
	//matchApi.DELETE("/:id", h.DeleteMatch)
	//
	//e.Logger.Fatal(e.Start(constant.AppPort))
}

func hello(w http.ResponseWriter, r *http.Request) {
	response.SetHTTPResponse(w, http.StatusOK, "success", nil)
	return
}
