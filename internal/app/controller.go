package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var validate *validator.Validate

type Controller struct {
	host    string
	service *Service
}

func init() {
	validate = validator.New()
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

func (controller *Controller) PostShortingJSON(c *gin.Context) {

	var jsonBody ShortingJSONBody

	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, controller.service.getHashByURL(jsonBody.URL))
	c.JSON(http.StatusCreated, gin.H{"result": shortenedURL})
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
