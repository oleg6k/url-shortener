package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oleg6k/url-shortener/internal/app"
	"github.com/oleg6k/url-shortener/internal/app/config"
	"github.com/oleg6k/url-shortener/internal/app/middleware"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()
	cfg := config.Load()

	repository, err := app.NewStorage("", cfg.FileStoragePath)
	if err != nil {
		sugar.Panicw(err.Error(), "event", "load repository")
	}

	service := app.NewService(repository)
	controller := app.NewController(cfg.BaseURL, service)

	r := gin.Default()

	r.Use(middleware.LoggerMiddleware(sugar))
	r.Use(middleware.GzipMiddleware())
	r.POST("/api/shorten", controller.PostShortingJSON)
	r.POST("/", controller.PostShorting)
	r.GET("/:shortUrl", controller.GetRedirectToOriginal)

	if err = r.Run(cfg.RunAddr); err != nil {
		sugar.Panicw(err.Error(), "event", "start server")
	}
}
