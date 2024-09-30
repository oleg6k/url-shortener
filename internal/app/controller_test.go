package app

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/oleg6k/url-shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fields struct {
	service *Service
	host    string
}
type want struct {
	code        int
	response    string
	contentType string
}

var cfg = config.Load()

func TestController_GetRedirectToOriginal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	tests := []struct {
		name   string
		want   want
		fields fields
	}{
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "invalid URL",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, err := NewStorage(cfg.Storage)
			assert.NoError(t, err)
			service := NewService(storage)
			controller := &Controller{
				host:    test.fields.host,
				service: service,
			}
			router.GET("/AcDbS", controller.GetRedirectToOriginal)
			req, err := http.NewRequest(http.MethodGet, "/AcDbS", nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, test.want.code, w.Code)
			assert.Contains(t, w.Body.String(), test.want.response)
			assert.Contains(t, w.Header().Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestController_PostShorting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	tests := []struct {
		name   string
		want   want
		fields fields
	}{
		{
			name: "positive test #1",
			want: want{
				code:        201,
				response:    "",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage, err := NewStorage(cfg.Storage)
			assert.NoError(t, err)
			service := NewService(storage)
			controller := &Controller{
				host:    test.fields.host,
				service: service,
			}
			router.POST("/", controller.PostShorting)
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://test.test")))
			req.Header.Set("Content-Type", "text/plain")
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, string(resBody))
			assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
