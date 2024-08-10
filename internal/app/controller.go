package app

import (
	"fmt"
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

func (controller *Controller) PostShorting(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusBadRequest)
		return
	}

	if request.Header.Get("Content-Type") != "text/plain" {
		http.Error(writer, "Invalid content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	originalURL := string(body)
	if originalURL == "" {
		http.Error(writer, "Empty URL provided", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(writer, "Invalid URL Format", http.StatusBadRequest)
		return
	}

	shortenedURL := fmt.Sprintf("%s/%s", controller.host, controller.service.getHashByUrl(originalURL))

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(shortenedURL))
}

func (controller *Controller) GetRedirectToOriginal(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusBadRequest)
		return
	}

	path := request.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) > 1 && parts[1] != "" {
		hash := parts[1]
		originalURL := controller.service.getUrlByHash(hash)
		if originalURL != "" {
			http.Redirect(writer, request, originalURL, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Error(writer, "Invalid URL", http.StatusBadRequest)
}
