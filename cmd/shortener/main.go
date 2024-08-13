package main

import (
	"github.com/oleg6k/url-shortener/config"
	"github.com/oleg6k/url-shortener/internal/app"
)

import (
	"github.com/gin-gonic/gin"
)

func main() {
	flags := ParseFlags()
	cfg := config.Load(flags.RunAddr, flags.BaseURL)

	service := app.NewService(make(map[string]string))
	controller := app.NewController(cfg.BaseURL, service)

	r := gin.Default()

	r.POST("/", controller.PostShorting)
	r.GET("/:shortUrl", controller.GetRedirectToOriginal)

	r.Run(cfg.RunAddr)
}
