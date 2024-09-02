package app

import (
	"bytes"
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

func TestController_GetRedirectToOriginal(t *testing.T) {

	tests := []struct {
		name   string
		want   want
		fields fields
	}{
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    "Invalid URL",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := &Controller{
				host:    test.fields.host,
				service: NewService(make(map[string]string)),
			}

			request := httptest.NewRequest(http.MethodGet, "/AcDbS", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			controller.GetRedirectToOriginal(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), test.want.response)
			assert.Contains(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestController_PostShorting(t *testing.T) {

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
			controller := &Controller{
				host:    test.fields.host,
				service: NewService(make(map[string]string)),
			}

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("https://test.test")))
			request.Header.Set("Content-Type", "text/plain")
			// создаём новый Recorder
			w := httptest.NewRecorder()
			controller.PostShorting(w, request)

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
