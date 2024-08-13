package main

import (
	"github.com/oleg6k/url-shortener/internal/app"
	"github.com/oleg6k/url-shortener/internal/app/config"
)

import (
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	service := app.NewService(make(map[string]string))
	controller := app.NewController(cfg.BaseURL, service)

	r := gin.Default()

	r.POST("/", controller.PostShorting)
	r.GET("/:shortUrl", controller.GetRedirectToOriginal)

	r.Run(cfg.RunAddr)
}
