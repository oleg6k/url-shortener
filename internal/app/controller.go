package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Controller struct {
	host    string
	service *Service
}

func NewController(host string, service *Service) *Controller {
	return &Controller{host: host, service: service}
}

func (controller *Controller) PostShorting(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	contentType := c.Request.Header.Get("Content-Type")
	mediaType := strings.TrimSpace(strings.Split(contentType, ";")[0])
	if mediaType != "text/plain" {
		c.String(http.StatusBadRequest, "invalid content type")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to read request body")
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		c.String(http.StatusBadRequest, "empty URL provided")
		return
	}

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid URL Format")
		return
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, controller.service.getHashByURL(originalURL))
	c.String(http.StatusCreated, shortenedURL)
}

func (controller *Controller) GetRedirectToOriginal(c *gin.Context) {
	c.Header("Content-Type", "text/plain")

	shortURL := c.Param("shortUrl")
	originalURL := controller.service.getURLByHash(shortURL)
	if originalURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, originalURL)
		return
	}

	c.String(http.StatusBadRequest, "invalid URL provided")
}
