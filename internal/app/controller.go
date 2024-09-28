package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/oleg6k/url-shortener/internal/app/types"
	"io"
	"net/http"
	"net/url"
	"slices"
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
	if !slices.Contains([]string{"application/x-gzip", "text/plain"}, mediaType) {
		c.AbortWithStatus(http.StatusBadRequest)
		c.String(http.StatusBadRequest, "Content-Type must be text/plain or application/x-gzip")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		c.String(http.StatusInternalServerError, "failed to read request body")
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		c.String(http.StatusBadRequest, "empty URL provided")
		return
	}

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		c.String(http.StatusBadRequest, "invalid URL Format")
		return
	}

	hash, err := controller.service.getHashByURL(originalURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, hash)
	c.String(http.StatusCreated, shortenedURL)
}

func (controller *Controller) PostShortingJSON(c *gin.Context) {
	var jsonBody types.ShortingJSONBody
	byteBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = json.Unmarshal(byteBody, &jsonBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = validate.Struct(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := controller.service.getHashByURL(jsonBody.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, hash)
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

func (controller *Controller) GetPing(c *gin.Context) {
	err := controller.service.Ping()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
