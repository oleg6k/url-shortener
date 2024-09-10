package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				return
			}
			defer gz.Close()
			c.Request.Body = io.NopCloser(gz)
		}

		c.Next()

		if strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			contentType := c.Writer.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") {
				c.Writer.Header().Set("Content-Encoding", "gzip")
				c.Writer.Header().Set("Vary", "Accept-Encoding")

				gz := gzip.NewWriter(c.Writer)
				defer gz.Close()

				c.Writer = &gzipWriter{Writer: gz, ResponseWriter: c.Writer}
			}
		}
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	io.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
