package main

import (
	"github.com/firzatullahd/blog-api/internal/config"
	"github.com/firzatullahd/blog-api/internal/delivery/http/handler"
	"github.com/firzatullahd/blog-api/internal/delivery/http/route"
	"github.com/firzatullahd/blog-api/internal/repository"
	"github.com/firzatullahd/blog-api/internal/usecase"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func main() {
	conf := config.Load()
	logger.Init()
	masterDB, replicaDB := config.InitializeDB(&conf.DB)

	repo := repository.NewRepository(masterDB, replicaDB)
	usecase := usecase.NewUsecase(conf, repo)
	handler := handler.NewHandler(usecase)
	route.Serve(conf, handler)
}
