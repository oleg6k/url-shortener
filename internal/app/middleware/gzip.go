package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"slices"
	"strings"
)

type Writer struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w gin.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (c *Writer) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *Writer) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *Writer) Close() error {
	return c.zw.Close()
}

type Reader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*Reader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		r:  r,
		zr: zr,
	}, nil
}

func (c Reader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *Reader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		acceptOptions := []string{"*/*", "application/json", "text/html", "application/x-gzip", "text/plain"}
		supportsGzipByContent := slices.Contains(acceptOptions, c.Request.Header.Get("Accept"))

		if supportsGzip && supportsGzipByContent {
			cw := newCompressWriter(c.Writer)
			c.Writer = cw
			defer cw.Close()
		}

		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			c.Request.Body = cr
			defer cr.Close()

		}
		c.Next()
	}
}
