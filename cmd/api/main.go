package main

import (
	"fmt"
	"github.com/firzatullahd/blog-api/internal/config"
	"github.com/firzatullahd/blog-api/internal/delivery/http/route"
	"github.com/firzatullahd/blog-api/internal/utils/logger"
)

func main() {
	conf := config.Load()
	logger.Init()
	fmt.Printf("conf %+v", conf)
	fmt.Println()
	_, _ = config.InitializeDB(&conf.DB)

	//repo := repository.NewRepository(masterDB, replicaDB)
	//usecase := usecase.NewUsecase(conf, repo)
	//handler := handler.NewHandler(usecase)
	route.Serve(conf, nil)
}
