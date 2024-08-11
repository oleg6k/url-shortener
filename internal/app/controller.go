package app

import (
	"errors"
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
	contentType := c.Request.Header.Get("Content-Type")
	mediaType := strings.TrimSpace(strings.Split(contentType, ";")[0])
	if mediaType != "text/plain" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid content type"))
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("failed to read request body"))
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("empty URL provided"))
		return
	}

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid URL Format"))
		return
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, controller.service.getHashByURL(originalURL))

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write([]byte(shortenedURL))
}

func (controller *Controller) GetRedirectToOriginal(c *gin.Context) {
	shortURL := c.Param("shortUrl")
	originalURL := controller.service.getURLByHash(shortURL)
	if originalURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, originalURL)
		return
	}
	c.AbortWithError(http.StatusBadRequest, errors.New("invalid URL provided"))
}
