package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware handles both Gzip decompression for incoming requests and compression for outgoing responses.
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Decompression: Check if the request is Gzip compressed
		if c.GetHeader("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				return
			}
			defer gz.Close()
			c.Request.Body = io.NopCloser(gz)
		}

		contentType := c.Writer.Header().Get("Content-Type")

		if (strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")) && strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			c.Writer.Header().Set("Vary", "Accept-Encoding")

			gz := gzip.NewWriter(c.Writer)
			defer gz.Close()

			gzWriter := &gzipWriter{Writer: gz, ResponseWriter: c.Writer}
			c.Writer = gzWriter

			c.Next()

			gzWriter.Flush()
		} else {
			c.Next()
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

func (w *gzipWriter) Flush() {
	if flusher, ok := w.Writer.(interface{ Flush() error }); ok {
		flusher.Flush()
	}
}
