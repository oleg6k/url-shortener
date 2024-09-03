package main

import (
	"github.com/oleg6k/url-shortener/internal/app"
)

import (
	"github.com/gin-gonic/gin"
)

func main() {
	service := app.NewService(make(map[string]string))
	controller := app.NewController("http://localhost:8080", service)
	r := gin.Default()

	r.POST("/", controller.PostShorting)
	r.GET("/:shortUrl", controller.GetRedirectToOriginal)

	r.Run(":8080")
}
