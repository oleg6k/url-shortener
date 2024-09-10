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
		if strings.Contains(c.GetHeader("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer gr.Close()
			c.Request.Body = io.NopCloser(gr)
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		acceptOptions := []string{"*/*", "application/json", "text/html", "application/x-gzip"}
		contentType := c.GetHeader("Accept")
		supportsGzipByContent := false
		for _, option := range acceptOptions {
			if strings.Contains(contentType, option) {
				supportsGzipByContent = true
				break
			}
		}

		if supportsGzip && supportsGzipByContent {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			cw := gzip.NewWriter(c.Writer)
			defer cw.Close()

			c.Writer = &gzipWriter{Writer: cw, ResponseWriter: c.Writer}
		}

		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	io.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}
