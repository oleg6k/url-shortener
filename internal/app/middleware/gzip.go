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

func (w *Writer) Write(p []byte) (int, error) {
	return w.zw.Write(p)
}

func (w *Writer) WriteHeader(statusCode int) {
	w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	w.ResponseWriter.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *Writer) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

func (w *Writer) Flush() {
	w.zw.Flush()
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *Writer) Close() error {
	return w.zw.Close()
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

func (r *Reader) Read(p []byte) (int, error) {
	return r.zr.Read(p)
}

func (r *Reader) Close() error {
	if err := r.zr.Close(); err != nil {
		return err
	}
	return r.r.Close()
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
